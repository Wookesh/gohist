package objects

import (
	"go/ast"
	"time"
)

type History struct {
	Data map[string]FunctionHistory
}

type FunctionHistory struct {
	History []HistoryElement
}

type HistoryElement struct {
	Commit *Commit
	Func   *ast.FuncDecl
}

type Commit struct {
	SHA       string
	Author    string
	Committer string
	Parent    []*Commit
	Timestamp time.Time
}
