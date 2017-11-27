package diff

import (
	"go/ast"
	"go/token"
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
	return ColorChange{color, node.Pos(), node.End()}
}

type Coloring []ColorChange
