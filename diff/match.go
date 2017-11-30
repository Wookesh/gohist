package diff

import (
	"go/ast"
	"math"
	"reflect"
	"sort"

	"github.com/Sirupsen/logrus"
)

type matchingElem struct {
	stmt  ast.Node
	score float64
	match ast.Node
}

type matchingElems []matchingElem

func (e matchingElems) Len() int           { return len(e) }
func (e matchingElems) Less(i, j int) bool { return e[i].score < e[j].score }
func (e matchingElems) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type matching struct {
	prev ast.Node
	next ast.Node
}

func matchNodes(a, b []ast.Node, callFunc string) []matching {
	matched := make(map[ast.Node]ast.Node)
	var matchingList matchingElems
	for _, aStmt := range a {
		for _, bStmt := range b {
			score := compare(aStmt, bStmt)
			logrus.Debugln(callFunc+":", reflect.TypeOf(aStmt), aStmt, score, reflect.TypeOf(bStmt), bStmt)
			if score > 1/math.Phi {
				matchingList = append(matchingList, matchingElem{aStmt, score, bStmt})
			}
		}
	}

	sort.Sort(matchingList)

	used := make(map[ast.Node]bool)
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

	var result []matching
	for _, aStmt := range a {
		bStmt := matched[aStmt]
		result = append(result, matching{prev: aStmt, next: bStmt})
	}
	return result
}

func matchStmts(a, b []ast.Stmt) []matching {
	return matchNodes(stmtsToNodes(a), stmtsToNodes(b), "matchStmts")
}

func matchExprs(a, b []ast.Expr) []matching {
	return matchNodes(exprToNodes(a), exprToNodes(b), "matchExprs")
}

func stmtsToNodes(l []ast.Stmt) (nodes []ast.Node) {
	for _, e := range l {
		nodes = append(nodes, e)
	}
	return
}

func exprToNodes(l []ast.Expr) (nodes []ast.Node) {
	for _, e := range l {
		nodes = append(nodes, e)
	}
	return
}

func matchFields(a, b []*ast.Field) []matching {
	return matchNodes(fieldToNodes(a), fieldToNodes(b), "matchFields")
}

func fieldToNodes(l []*ast.Field) (nodes []ast.Node) {
	for _, e := range l {
		nodes = append(nodes, e)
	}
	return
}

func matchIdents(a, b []*ast.Ident) []matching {
	return matchNodes(identToNodes(a), identToNodes(b), "matchIdents")
}

func identToNodes(l []*ast.Ident) (nodes []ast.Node) {
	for _, e := range l {
		nodes = append(nodes, e)
	}
	return
}
