package diff

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/Sirupsen/logrus"
	"github.com/wookesh/gohist/util"
)

type Color int

const (
	ColorSame Color = iota
	ColorNew
	ColorRemoved
	ColorSimilar
)

type Mode int

const (
	ModeNew Mode = iota
	ModeOld
)

func (m Mode) String() string {
	switch m {
	case ModeNew:
		return "ModeNew"
	case ModeOld:
		return "ModeOld"
	default:
		return ""
	}
}

func (m Mode) ToColor() Color {
	switch m {
	case ModeNew:
		return ColorNew
	case ModeOld:
		return ColorRemoved
	default:
		return ColorSame
	}
}

type ColorChange struct {
	Color Color
	Pos   token.Pos
	End   token.Pos
}

func NewColorChange(color Color, node ast.Node) ColorChange {
	logrus.Debugln("NewColorChange:", color, node, node.Pos(), node.End()-1)
	return ColorChange{color, node.Pos(), node.End() - 1}
}

type Coloring []ColorChange

func colorMatches(matching []matching, mode Mode, callFunc string) (coloring Coloring) {
	for _, match := range matching {
		if match.next == nil {
			logrus.Debugln(callFunc, "unmatched:", match.prev, reflect.TypeOf(match.prev))
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			nodeColoring := diff(match.prev, match.next, mode)
			if len(nodeColoring) == 0 && match.orderChanged {
				nodeColoring = append(nodeColoring, NewColorChange(ColorSimilar, match.prev))
			}
			coloring = append(coloring, nodeColoring...)

		}
	}
	return
}

func colorList(a, b []ast.Node, mode Mode, callFunc string) (coloring Coloring) {
	min := util.IntMin(len(a), len(b))
	for i := 0; i < min; i++ {
		coloring = append(coloring, diff(a[i], b[i], mode)...)
	}
	if len(a) > min {
		for _, aNode := range a[min:] {
			coloring = append(coloring, NewColorChange(mode.ToColor(), aNode))
		}
	}
	return
}
