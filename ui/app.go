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
	history *objects.History
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

type Link struct {
	Name string
	Len  int
}

type Links []Link

func (l Links) Len() int           { return len(l) }
func (l Links) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }
func (l Links) Less(i, j int) bool { return l[i].Name < l[j].Name }

func (h *handler) List(c echo.Context) error {
	var items Links
	for i, fHistory := range h.history.Data {
		items = append(items, Link{i, len(fHistory.History)})
	}
	sort.Sort(items)
	return c.Render(http.StatusOK, "list.html", items)
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

func Run(history *objects.History) {
	handler := handler{history: history}

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
	}

	t := &Template{
		templates: template.Must(template.New("sites").Funcs(funcMap).ParseGlob("ui/views/*.html")),
	}
	e := echo.New()
	e.Renderer = t
	e.GET("/", handler.List)
	e.Static("/static", "ui/static")
	e.GET("/:name/", handler.Get)
	e.GET("/:path/:name/", handler.Get)
	e.Logger.Fatal(e.Start(":8000"))
}

func color(s string, coloring diff.Coloring, offset int) template.HTML {
	if len(coloring) == 0 {
		return template.HTML(s)
	}
	logrus.Debugln("color:")
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
		result += string(s[i])
	}
	if hasColoring {
		result += `</span>`
	}
	//fmt.Println(result)
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
