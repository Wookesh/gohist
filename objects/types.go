package objects

import (
	"fmt"
	"go/ast"
	"html/template"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type History struct {
	Data            map[string]*FunctionHistory
	CommitsAnalyzed int
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
		changes += len(history.History) - 1
		if len(history.History) == 1 {
			neverChanged++
		}
		if history.Deleted {
			removed++
		}
		if len(history.History) > mostChangedCount {
			mostChanged = name
			mostChangedCount = len(history.History)
		}
	}
	stats["Analyzed commits"] = h.CommitsAnalyzed
	stats["Changes per commit"] = changes / h.CommitsAnalyzed
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
		changesCount[len(fHistory.History)] += 1
		for _, commit := range fHistory.History {
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

	logrus.Infoln(changesCount)

	return charts
}

type FunctionHistory struct {
	History         []*HistoryElement
	LifeTime        int
	FirstAppearance time.Time
	LastAppearance  time.Time
	Deleted         bool
}

type HistoryElement struct {
	Commit *object.Commit
	Func   *ast.FuncDecl
	Text   string
	Offset int
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
