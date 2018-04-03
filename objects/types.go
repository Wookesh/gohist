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

	"github.com/wookesh/gohist/diff"
	"github.com/wookesh/gohist/util"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type History struct {
	Data            map[string]*FunctionHistory
	CommitsAnalyzed int32

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

func (h *History) CheckForDeleted(commit *object.Commit) {
	for _, fh := range h.Data {
		fh.Delete(commit)
	}
}

func NewHistory() *History {
	return &History{
		Data: make(map[string]*FunctionHistory),
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
		changes += len(history.Elements) - 1
		if len(history.Elements) == 1 {
			neverChanged++
		}
		if history.Deleted {
			removed++
		}
		if len(history.Elements) > mostChangedCount {
			mostChanged = name
			mostChangedCount = len(history.Elements)
		}
	}
	stats["Analyzed commits"] = h.CommitsAnalyzed
	stats["Changes per commit"] = changes / int(h.CommitsAnalyzed)
	stats["Changes per function"] = float64(changes) / float64(len(h.Data))
	stats["Never changed"] = neverChanged
	stats["Functions"] = len(h.Data)
	stats["Most changed"] = fmt.Sprintf("%v [%v]", mostChanged, mostChangedCount)
	stats["Removed"] = removed
	return stats
}

type ChartData struct {
	X     template.JS
	YAxis string
	Y     template.JS
	Type  string
	Name  string
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

func (h *History) ChartsData() map[string]ChartData {
	charts := make(map[string]ChartData)

	changesCount := make(map[int]int)
	changedPerDate := make(map[Date]int)
	for _, fHistory := range h.Data {
		changesCount[len(fHistory.Elements)] += 1
		for _, commit := range fHistory.Elements {
			var date Date
			date.Year, date.Month, date.Day = commit.Commit.Author.When.Date()
			changedPerDate[date] += 1
		}
	}
	xAxis, yAxis := toStrings(changesCount)
	charts["function_changed_count"] = ChartData{X: template.JS(xAxis), Y: template.JS(yAxis), YAxis: "functions count",
		Name: "function changed count"}

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

	charts["functions_changed_per_day"] = ChartData{X: template.JS(strings.Join(xAxis2List, ",")),
		Y: template.JS(strings.Join(yAxis2List, ",")), YAxis: "functions changed", Type: "timeseries",
		Name: "functions changed per day",
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
	parentMapping map[string][]string
	m             sync.Mutex
}

func NewFunctionHistory(id string) *FunctionHistory {
	return &FunctionHistory{
		Elements:      make(map[string]*HistoryElement),
		parentMapping: make(map[string][]string),
		ID:            id,
	}
}

func (fh *FunctionHistory) AddElement(decl *ast.FuncDecl, commit *object.Commit, body []byte) {
	fh.m.Lock()
	defer fh.m.Unlock()

	sha := commit.Hash.String()
	fh.LifeTime++

	parents := make(map[string]*HistoryElement)
	// physical parent
	anyDifferent := false
	var parentMapping []string
	for _, parent := range commit.ParentHashes {
		parentSHA := parent.String()
		mapped, ok := fh.parentMapping[parentSHA]
		if !ok {
			continue
		}
		// logical parent
		for _, parent := range mapped {
			parent, ok := fh.Elements[parent]
			if !ok {
				continue
			}
			if diff.IsSame(parent.Func, decl) {
				parentMapping = append(parentMapping, parent.Commit.Hash.String())
			} else {
				anyDifferent = true
				parents[parentSHA] = parent
			}
		}
	}
	if !anyDifferent && len(fh.Elements) > 0 {
		fh.parentMapping[sha] = parentMapping
		return
	}
	element := &HistoryElement{
		Func:     decl,
		Commit:   commit,
		Parent:   parents,
		Children: make(map[string]*HistoryElement),
		Text:     string(body[decl.Pos()-1 : decl.End()-1]),
		Offset:   int(decl.Pos()),
	}

	for _, parent := range parents {
		parent.Children[sha] = element
	}
	fh.Elements[sha] = element
	fh.parentMapping[sha] = []string{sha}
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
		for _, parent := range mapped {
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
	}

	for _, parent := range parents {
		parent.Children[sha] = element
	}
	fh.Elements[sha] = element
	fh.parentMapping[sha] = []string{sha}
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

type HistoryElement struct {
	Commit *object.Commit
	Func   *ast.FuncDecl
	Text   string
	Offset int

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
