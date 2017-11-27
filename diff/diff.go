package diff

import (
	"go/ast"
	"log"
	"reflect"
	"sort"
)

func Diff(a, b ast.Node, mode Mode) Coloring {
	if mode == ModeNew && a == nil {
		return Coloring{NewColorChange(ColorNew, b)}
	}
	if mode == ModeOld && b == nil {
		return Coloring{NewColorChange(ColorRemoved, a)}
	}
	return diff(a, b, mode)
}

func diff(a, b ast.Node, mode Mode) Coloring {
	var coloring Coloring
	switch t := a.(type) {
	case ast.Decl:
		coloring = diffDecl(t, b, mode)
	case ast.Stmt:
		coloring = diffStmt(t, b, mode)
		//case ast.Expr:
		//	diffExpr(t, b, mode)
		//case ast.Decl:
		//	diffDecl(t, b, mode)
	default:
		coloring = Coloring{NewColorChange(ColorSame, a)}
	}

	return coloring
}

func diffDecl(a ast.Decl, b ast.Node, mode Mode) Coloring {
	_, ok := b.(ast.Decl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	switch t := a.(type) {
	case *ast.FuncDecl:
		return diffFuncDecl(t, b, mode)
	}
	return nil
}

func diffFuncDecl(a *ast.FuncDecl, bNode ast.Node, mode Mode) Coloring {
	b, ok := bNode.(*ast.FuncDecl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a.Body, b.Body, mode)
}

func diffStmt(a ast.Stmt, b ast.Node, mode Mode) Coloring {
	_, ok := b.(ast.Stmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	switch t := a.(type) {
	case *ast.BlockStmt:
		return diffBlockStmt(t, b, mode)
	}
	return nil
}

func diffBlockStmt(a *ast.BlockStmt, bNode ast.Node, mode Mode) Coloring {
	b, ok := bNode.(*ast.BlockStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	var coloring Coloring
	for aStmt, bStmt := range matchStmts(a.List, b.List) {
		if bStmt == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), aStmt))
		}
		coloring = append(coloring, diff(aStmt, bStmt, mode)...)
	}
	return coloring
}

type matchingElem struct {
	stmt  ast.Stmt
	score float64
	match ast.Stmt
}

type matchingElems []matchingElem

func (e matchingElems) Len() int {
	return len(e)
}

func (e matchingElems) Less(i, j int) bool {
	return e[i].score < e[j].score
}

func (e matchingElems) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func matchStmts(a, b []ast.Stmt) map[ast.Stmt]ast.Stmt {
	matched := make(map[ast.Stmt]ast.Stmt)
	var matchingList matchingElems
	for _, aStmt := range a {
		for _, bStmt := range b {
			score := compare(aStmt, bStmt)
			if score > 0.0 {
				log.Println("matchStmts:", reflect.TypeOf(aStmt), aStmt, score, reflect.TypeOf(bStmt), bStmt)
				matchingList = append(matchingList, matchingElem{aStmt, score, bStmt})
			}
		}
	}

	sort.Sort(matchingList)

	used := make(map[ast.Stmt]bool)
	for _, elem := range matchingList {
		if _, ok := used[elem.match]; ok {
			continue
		}
		if _, ok := matched[elem.stmt]; ok {
			continue
		}
		used[elem.match] = true
		matched[elem.stmt] = elem.match
	}

	for _, aStmt := range a {
		if _, ok := matched[aStmt]; !ok {
			matched[aStmt] = nil
		}
	}
	return matched
}
