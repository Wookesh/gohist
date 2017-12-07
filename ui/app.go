package ui

import (
	"html/template"
	"io"
	"net/http"
	"sort"
	"strconv"

	"github.com/Sirupsen/logrus"
	"github.com/labstack/echo"
	"github.com/wookesh/gohist/diff"
	"github.com/wookesh/gohist/objects"
)

type handler struct {
	history  *objects.History
	repoName string
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Link struct {
	Name  string
	Len   int
	Total int
}

type ListData struct {
	RepoName string
	Links    Links
}
type Links []Link

func (l Links) Len() int           { return len(l) }
func (l Links) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Links) Less(i, j int) bool { return l[i].Name < l[j].Name }

func (h *handler) List(c echo.Context) error {
	onlyChangedStr := c.QueryParam("only_changed")
	onlyChanged, err := strconv.ParseBool(onlyChangedStr)
	if err != nil {
		onlyChanged = false
	}
	listData := &ListData{RepoName: h.repoName}
	for i, fHistory := range h.history.Data {
		if !onlyChanged || (onlyChanged && len(fHistory.History) > 1) {
			listData.Links = append(listData.Links,
				Link{
					Name:  i,
					Len:   len(fHistory.History),
					Total: fHistory.LifeTime,
				})
		}
	}
	sort.Sort(listData.Links)
	return c.Render(http.StatusOK, "list.html", listData)
}

type DiffView struct {
	Name      string
	History   *objects.FunctionHistory
	LeftDiff  diff.Coloring
	RightDiff diff.Coloring
}

func (h *handler) Get(c echo.Context) error {
	funcName := c.Param("name")
	pack := c.Param("path")
	if pack != "" {
		funcName = pack + "/" + funcName
	}
	f, ok := h.history.Data[funcName]
	if !ok {
		c.HTML(http.StatusNotFound, "NOT FOUND")
	}
	var pos int64
	pos, err := strconv.ParseInt(c.QueryParam("pos"), 10, 32)
	if err != nil {
		pos = 0
	}
	var left, right diff.Coloring
	if pos > 0 {
		left = diff.Diff(f.History[pos-1].Func, f.History[pos].Func, diff.ModeOld)
		right = diff.Diff(f.History[pos].Func, f.History[pos-1].Func, diff.ModeNew)
	} else {
		right = diff.Diff(nil, f.History[pos].Func, diff.ModeNew)
	}
	diffView := &DiffView{Name: funcName, History: f, LeftDiff: left, RightDiff: right}
	data := map[string]interface{}{"pos": pos, "diffView": diffView}
	return c.Render(http.StatusOK, "diff.html", data)
}

func Run(history *objects.History, repoName string) {
	handler := handler{history: history, repoName: repoName}

	funcMap := template.FuncMap{
		"next": func(i int64) int64 {
			return i + 1
		},
		"prev": func(i int64) int64 {
			return i - 1
		},
		"prev_int": func(i int) int {
			return i - 1
		},
		"color": color,
		"modifications": func(a, b int) string {
			if b == 0 {
				return "dark"
			}
			stability := 1.0 - float64(a)/float64(b)
			if stability >= 0.8 {
				return "success"
			} else if stability >= 0.5 {
				return "warning"
			} else {
				return "danger"
			}
		},
	}

	t := &Template{
		templates: template.Must(template.New("sites").Funcs(funcMap).ParseGlob("ui/views/*.html")),
	}
	e := echo.New()
	e.HideBanner = true
	e.Renderer = t

	e.GET("/", handler.List)
	e.GET("/:name/", handler.Get)
	e.GET("/:path/:name/", handler.Get)
	e.Static("/static", "ui/static")

	logrus.Infoln("GoHist:", "started web server")

	if err := e.Start(":8000"); err != nil {
		logrus.Fatalln(err)
	}
}

func color(s string, coloring diff.Coloring, offset int) template.HTML {
	if len(coloring) == 0 {
		return template.HTML(s)
	}
	logrus.Debugln("color:", coloring, offset)
	current := 0
	var hasColoring bool
	var result string
	logrus.Debugln("color:", "next coloring:", current, coloring[current])
	for i := 0; i < len(s); i++ {
		if current < len(coloring) {
			if !hasColoring && int(coloring[current].Pos) <= i+offset {
				logrus.Debugln("color:", "changing color:", toColor(coloring[current].Color), i+offset)
				hasColoring = true
				result += `<span style="color: ` + toColor(coloring[current].Color) + `;">`
			}

			if hasColoring && int(coloring[current].End) < i+offset {
				logrus.Debugln("color:", "removing color:", i+offset)
				result += `</span>`
				if current < len(coloring) {
					current++
					logrus.Debugln("color:", "next coloring:", current)
				}
				if current < len(coloring) && int(coloring[current].Pos) <= i+offset {
					logrus.Debugln("color:", "changing color:", toColor(coloring[current].Color), i+offset)
					result += `<span style="color: ` + toColor(coloring[current].Color) + `;">`
				} else {
					hasColoring = false
				}
			}
		}
		if s[i] == '<' {
			result += `<span>` + string(s[i]) + `</span>` // TODO: I dunno how to frontend, find better solution later
		} else {
			result += string(s[i])
		}
	}
	if hasColoring {
		result += `</span>`
	}
	return template.HTML(result)
}

func toColor(c diff.Color) string {
	switch c {
	case diff.ColorSame:
		return "white"
	case diff.ColorNew:
		return "green"
	case diff.ColorRemoved:
		return "red"
	case diff.ColorSimilar:
		return "lightblue"
	default:
		return "white"
	}
}
