package objects

import (
	"go/ast"

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
	History []*HistoryElement
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
