package diff

import (
	"go/ast"
	"go/token"
	"reflect"
	"strings"

	"github.com/sirupsen/logrus"
)

type context struct {
	a nodeContext
	b nodeContext
}

type nodeContext struct {
	vars       vars
	globalVars vars
}

type vars map[string][]token.Pos

type simpleColoring struct {
	Color Color
	Data  string
}

func LCS(a, b string, offset int, mode Mode) (coloring Coloring) {
	aList := strings.Split(a, "\n")
	bList := strings.Split(b, "\n")
	c := lcs(aList, bList)
	for _, simpleColor := range printLCS(c, aList, bList, len(aList), len(bList), mode) {
		logrus.Infoln(simpleColor.Color, simpleColor.Data)
		coloring = append(coloring, ColorChange{Color: simpleColor.Color, Pos: token.Pos(offset), End: token.Pos(offset + len(simpleColor.Data))})
		offset += len(simpleColor.Data) + 1
	}
	return
}

func printLCS(c [][]int, x, y []string, i, j int, mode Mode) (coloring []simpleColoring) {
	logrus.Infoln(i, j, len(x), len(y))
	if i > 0 && j > 0 && x[i-1] == y[j-1] {
		coloring = printLCS(c, x, y, i-1, j-1, mode)
		coloring = append(coloring, simpleColoring{Color: ColorSame, Data: x[i-1]})
		logrus.Infoln(" " + x[i-1])
	} else {
		if j > 0 && (i == 0 || c[i][j-1] >= c[i-1][j]) {
			coloring = printLCS(c, x, y, i, j-1, mode)
			if mode == ModeNew {
				coloring = append(coloring, simpleColoring{Color: ColorNew, Data: y[j-1]})
				logrus.Infoln("+" + y[j-1])
			}
		} else if i > 0 && (j == 0 || c[i][j-1] < c[i-1][j]) {
			coloring = printLCS(c, x, y, i-1, j, mode)
			if mode == ModeOld {
				coloring = append(coloring, simpleColoring{Color: ColorRemoved, Data: x[i-1]})
				logrus.Infoln("-" + x[i-1])
			}
		}
	}
	return coloring
}

func lcs(x, y []string) [][]int {
	m := len(x)
	n := len(y)
	c := make([][]int, m+1)
	for i := range c {
		c[i] = make([]int, n+1)
	}
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if x[i] == y[j] {
				c[i+1][j+1] = c[i][j] + 1
			} else {
				c[i+1][j+1] = max(c[i+1][j], c[i][j+1])
			}
		}
	}
	return c
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Diff(a, b ast.Node, mode Mode) Coloring {
	logrus.Debugln("Diff:", mode, "\n")
	if mode == ModeNew && a == nil {
		return Coloring{NewColorChange(mode.ToColor(), b)}
	}
	if mode == ModeOld && b == nil {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a, b, mode)
}

func diff(aNode, b ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diff:", aNode, b)
	if aNode == nil {
		return
	}
	switch a := aNode.(type) {
	case ast.Decl:
		coloring = diffDecl(a, b, mode)
	case ast.Expr:
		coloring = diffExpr(a, b, mode)
	case ast.Stmt:
		coloring = diffStmt(a, b, mode)
	// non interface nodes:
	case *ast.Field:
		coloring = diffField(a, b, mode)
	case *ast.FieldList:
		coloring = diffFieldList(a, b, mode)
	case *ast.ValueSpec:
		coloring = diffValueSpec(a, b, mode)
	default:
		logrus.Errorln("diff:", "not implemented case", reflect.TypeOf(a))
		coloring = Coloring{NewColorChange(mode.ToColor(), a)}
	}

	return
}

func diffFieldList(a *ast.FieldList, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.FieldList)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a == nil {
		return
	}
	if b == nil {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	for _, match := range matchFields(a.List, b.List) {
		if match.next == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}

	return
}

func diffField(a *ast.Field, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.Field)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	coloring = append(coloring, colorMatches(matchIdents(a.Names, b.Names), mode, "diffField")...)
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	return
}

func diffValueSpec(a *ast.ValueSpec, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.ValueSpec)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	coloring = append(coloring, colorMatches(matchIdents(a.Names, b.Names), mode, "diffValueSpec")...)
	coloring = append(coloring, colorMatches(matchExprs(a.Values, b.Values), mode, "diffValueSpec")...)
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	return
}
