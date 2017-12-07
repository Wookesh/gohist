package objects

import (
	"go/ast"
	"time"

	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type History struct {
	Data map[string]*FunctionHistory
}

func NewHistory() *History {
	return &History{
		Data: make(map[string]*FunctionHistory),
	}
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
