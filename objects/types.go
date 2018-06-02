package objects

import (
	"fmt"
	"go/ast"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/gohist/diff"
	"github.com/wookesh/gohist/util"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type History struct {
	Data            map[string]*FunctionHistory
	CommitsAnalyzed int32
	MaxChanged      int32
	CountPerCommit  map[time.Time]int

	m sync.Mutex
}

func (h *History) Get(funcID string) *FunctionHistory {
	h.m.Lock()
	defer h.m.Unlock()
	funcHistory, ok := h.Data[funcID]
	if !ok {
		funcHistory = NewFunctionHistory(funcID)
		h.Data[funcID] = funcHistory
	}
	return funcHistory
}

func (h *History) Mark(sha time.Time, count int) {
	h.m.Lock()
	h.CountPerCommit[sha] = count
	h.m.Unlock()
}

func (h *History) CheckForDeleted(commit *object.Commit) {
	h.m.Lock()
	defer h.m.Unlock()
	for _, fh := range h.Data {
		fh.Delete(commit)
	}
}

func NewHistory() *History {
	return &History{
		Data:           make(map[string]*FunctionHistory),
		CountPerCommit: make(map[time.Time]int),
	}
}

func (h *History) Stats() map[string]interface{} {
	stats := make(map[string]interface{})
	changes := 0
	neverChanged := 0
	mostChangedCount := 0
	removed := 0
	var mostChanged string
	for name, history := range h.Data {
		versions := history.VersionsCount()
		changes += versions - 1
		if versions == 1 {
			neverChanged++
		}
		if history.Deleted {
			removed++
		}
		if versions > mostChangedCount {
			mostChanged = name
			mostChangedCount = versions
		}
	}
	stats["Analyzed commits"] = h.CommitsAnalyzed
	stats["Changes per commit"] = float64(changes) / float64(h.CommitsAnalyzed)
	stats["Changes per function"] = float64(changes) / float64(len(h.Data))
	stats["Max changes in commit"] = h.MaxChanged
	stats["Never changed"] = neverChanged
	stats["Functions"] = len(h.Data)
	stats["Most changed"] = fmt.Sprintf("%v [%v]", mostChanged, mostChangedCount)
	stats["Removed"] = removed
	logrus.Infof("%v,%v,%v,%v,%v,%v,%v,%v",
		stats["Analyzed commits"],
		stats["Changes per commit"],
		stats["Changes per function"],
		stats["Max changes in commit"],
		stats["Never changed"],
		stats["Functions"],
		mostChangedCount,
		stats["Removed"],
	)

	return stats
}

type PieRow struct {
	Name  string
	Value int
}

func PieRowsFromMap(m map[string]int) (rows []PieRow) {
	rows = append(rows, PieRow{"stable", m["stable"]})
	rows = append(rows, PieRow{"modified", m["modified"]})
	rows = append(rows, PieRow{"active", m["active"]})
	return
}

type ChartData struct {
	X       template.JS
	YAxis   string
	Y       template.JS
	Type    string
	Name    string
	PieData []PieRow
}

type Date struct {
	Year  int
	Month time.Month
	Day   int
}

func (d Date) String() string {
	s := time.Date(d.Year, d.Month, d.Day, 0, 0, 0, 0, time.Local).Format("'2006-01-02'")
	return s
}

type Dates []Date

func (d Dates) Len() int { return len(d) }

func (d Dates) Less(i, j int) bool {
	a := d[i]
	b := d[j]
	return a.Year < b.Year || (a.Year == b.Year && (a.Month < b.Month || (a.Month == b.Month && a.Day < b.Day)))
}

func (d Dates) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

type DateCount struct {
	Date  time.Time
	Count int
}

type DateCounts []DateCount

func (d DateCounts) Len() int {
	return len(d)
}

func (d DateCounts) Less(i, j int) bool {
	return d[i].Date.Before(d[j].Date)
}

func (d DateCounts) Swap(i, j int) {
	d[i], d[j] = d[j], d[i]
}

func ToStabilityGroup(stability float64) string {
	if stability >= 0.8 {
		return "stable"
	}
	if stability >= 0.5 {
		return "modified"
	}
	return "active"
}

func (h *History) ChartsData() map[string]ChartData {
	charts := make(map[string]ChartData)

	changesCount := make(map[int]int)
	changedPerDate := make(map[Date]int)
	stabilityVersions := map[string]int{"stable": 0, "modified": 0, "active": 0}
	for _, fHistory := range h.Data {
		changesCount[fHistory.VersionsCount()] += 1
		stability := 1.0 - float64(fHistory.VersionsCount())/float64(fHistory.LifeTime)
		stabilityVersions[ToStabilityGroup(stability)] += 1
		for _, commit := range fHistory.Elements {
			var date Date
			date.Year, date.Month, date.Day = commit.Commit.Author.When.Date()
			changedPerDate[date] += 1
		}
	}
	xAxis, yAxis := toStrings(changesCount)
	charts["function_changed_count"] = ChartData{
		X:     template.JS(xAxis),
		Y:     template.JS(yAxis),
		YAxis: "functions count",
		Name:  "function changed count",
		Type:  "common",
	}

	charts["stability_chart"] = ChartData{
		Type:    "pie",
		PieData: PieRowsFromMap(stabilityVersions),
		Name:    "stability",
	}

	var ordered Dates
	for k := range changedPerDate {
		ordered = append(ordered, k)
	}
	sort.Sort(ordered)
	var xAxis2List, yAxis2List []string
	for _, date := range ordered {
		yAxis2List = append(yAxis2List, strconv.FormatInt(int64(changedPerDate[date]), 10))
		xAxis2List = append(xAxis2List, date.String())
	}

	charts["functions_changed_per_day"] = ChartData{
		X:     template.JS(strings.Join(xAxis2List, ",")),
		Y:     template.JS(strings.Join(yAxis2List, ",")),
		YAxis: "functions changed",
		Type:  "timeseries",
		Name:  "functions changed per day",
	}

	var dateCounts DateCounts
	for date, count := range h.CountPerCommit {
		dateCounts = append(dateCounts, DateCount{date, count})
	}

	sort.Sort(dateCounts)

	var dates, counts []string
	for _, d := range dateCounts {
		dates = append(dates, "'"+d.Date.Format("2006-01-02T15:04:05")+"'")
		counts = append(counts, strconv.FormatInt(int64(d.Count), 10))
	}

	charts["functions_count_in_time"] = ChartData{
		Y:     template.JS(strings.Join(counts, ",")),
		X:     template.JS(strings.Join(dates, ",")),
		YAxis: "functions count",
		Type:  "datetimeseries",
		Name:  "functions count in time",
	}

	return charts
}

type FunctionHistory struct {
	History         []*HistoryElement
	LifeTime        int
	FirstAppearance time.Time
	LastAppearance  time.Time
	Deleted         bool

	ID            string
	Elements      map[string]*HistoryElement
	First, Last   *HistoryElement
	parentMapping map[string]map[string]bool
	m             sync.Mutex
}

func NewFunctionHistory(id string) *FunctionHistory {
	return &FunctionHistory{
		Elements:      make(map[string]*HistoryElement),
		parentMapping: make(map[string]map[string]bool),
		ID:            id,
	}
}

func (fh *FunctionHistory) AddElement(decl *ast.FuncDecl, commit *object.Commit, body []byte) bool {
	fh.m.Lock()
	defer fh.m.Unlock()

	sha := commit.Hash.String()
	fh.LifeTime++

	parents := make(map[string]*HistoryElement)
	// physical parent
	anyDifferent := false
	anySame := false
	parentMapping := make(map[string]bool)
	for _, parent := range commit.ParentHashes {
		parentSHA := parent.String()
		mapped, ok := fh.parentMapping[parentSHA]
		if !ok {
			continue
		}
		// logical parent
		for parent := range mapped {
			parentSHA = parent
			parent, ok := fh.Elements[parent]
			if !ok {
				continue
			}
			parents[parentSHA] = parent
			if diff.IsSame(parent.Func, decl) {
				anySame = true
				parentMapping[parent.Commit.Hash.String()] = true
			} else {
				anyDifferent = true
			}
		}
	}
	if !anyDifferent && len(fh.Elements) > 0 {
		fh.parentMapping[sha] = parentMapping
		return false
	}
	element := &HistoryElement{
		Func:     decl,
		Commit:   commit,
		Parent:   parents,
		Children: make(map[string]*HistoryElement),
		Text:     string(body[decl.Pos()-1 : decl.End()-1]),
		Offset:   int(decl.Pos()),
		New:      !anySame,
	}

	for _, parent := range parents {
		parent.Children[sha] = element
	}
	fh.Elements[sha] = element
	fh.parentMapping[sha] = map[string]bool{sha: true}
	if fh.Deleted {
		fh.Deleted = false
	}
	return !anySame
}

func (fh *FunctionHistory) Delete(commit *object.Commit) {
	fh.m.Lock()
	defer fh.m.Unlock()

	_, ok := fh.parentMapping[commit.Hash.String()]
	if ok {
		return
	}

	sha := commit.Hash.String()

	parents := make(map[string]*HistoryElement)
	var anyNotDeleted bool
	// physical parent
	for _, parent := range commit.ParentHashes {
		parentSHA := parent.String()
		mapped, ok := fh.parentMapping[parentSHA]
		if !ok {
			continue
		}
		// logical parent
		for parent := range mapped {
			parentSHA = parent
			parent, ok := fh.Elements[parent]
			if !ok {
				continue
			}
			if parent.Func != nil {
				anyNotDeleted = true
			}
			parents[parentSHA] = parent
		}
	}
	if !anyNotDeleted {
		return
	}
	element := &HistoryElement{
		Func:     nil,
		Commit:   commit,
		Parent:   parents,
		Children: make(map[string]*HistoryElement),
		New:      false,
	}

	for _, parent := range parents {
		parent.Children[sha] = element
	}
	fh.Elements[sha] = element
	fh.parentMapping[sha] = map[string]bool{sha: true}
	fh.Deleted = true
}

func (fh *FunctionHistory) PostProcess() {
	fh.m.Lock()
	defer fh.m.Unlock()
	for _, elem := range fh.Elements {
		thisTime := util.Earlier(elem.Commit.Author.When, elem.Commit.Committer.When)
		if fh.First == nil {
			fh.First = elem
		} else {
			if thisTime.Before(util.Earlier(fh.First.Commit.Author.When, fh.First.Commit.Committer.When)) {
				fh.First = elem
			}
		}
		if fh.Last == nil {
			fh.Last = elem
		} else {
			if util.Earlier(fh.Last.Commit.Author.When, fh.Last.Commit.Committer.When).Before(thisTime) {
				fh.Last = elem
			}
		}
	}
}

func (fh *FunctionHistory) VersionsCount() int {
	versions := 0
	for _, elem := range fh.Elements {
		if elem.New {
			versions++
		}
	}
	return versions
}

type HistoryElement struct {
	Commit *object.Commit
	Func   *ast.FuncDecl
	Text   string
	Offset int
	New    bool

	Parent   map[string]*HistoryElement
	Children map[string]*HistoryElement
}

type Variable struct {
	Name *ast.Ident
	Type ast.Expr
	Expr ast.Expr
}

func toStrings(m map[int]int) (string, string) {
	var a, b []string
	for k, v := range m {
		a = append(a, strconv.FormatInt(int64(k), 10))
		b = append(b, strconv.FormatInt(int64(v), 10))
	}
	return strings.Join(a, ","), strings.Join(b, ",")
}
