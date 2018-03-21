package objects

import (
	"fmt"
	"go/ast"
	"html/template"
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
}

func (h *History) ChartsData() map[string]ChartData {
	charts := make(map[string]ChartData)

	changesCount := make(map[int]int)
	for _, fHistory := range h.Data {
		changesCount[len(fHistory.History)] += 1
	}
	xAxis, yAxis := toStrings(changesCount)
	charts["changes_to_commits"] = ChartData{X: template.JS(xAxis), Y: template.JS(yAxis), YAxis: "functions count"}

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
