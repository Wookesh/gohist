package diff

import (
	"go/ast"
	"reflect"
	"sort"

	"github.com/sirupsen/logrus"
)

type matchingElem struct {
	node  *nodeInfo
	score float64
	match *nodeInfo
}

type nodeInfo struct {
	node ast.Node
	pos  int
}

type matchingElems []matchingElem

func (e matchingElems) Len() int           { return len(e) }
func (e matchingElems) Less(i, j int) bool { return e[i].score < e[j].score }
func (e matchingElems) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

type matching struct {
	prev            ast.Node
	next            ast.Node
	positionChanged bool
	orderChanged    bool
}

func matchNodes(a, b []ast.Node, callFunc string) []matching {
	var matchingList matchingElems
	var aNodesInfo, bNodeInfo []*nodeInfo
	for i, node := range a {
		aNodesInfo = append(aNodesInfo, &nodeInfo{node, i})
	}
	for i, node := range b {
		bNodeInfo = append(bNodeInfo, &nodeInfo{node, i})
	}
	for _, aNodeInfo := range aNodesInfo {
		for _, bNodeInfo := range bNodeInfo {
			score := compare(aNodeInfo.node, bNodeInfo.node)
			logrus.Debugln(callFunc+":", reflect.TypeOf(aNodeInfo.node), aNodeInfo.node, score,
				reflect.TypeOf(bNodeInfo.node), bNodeInfo.node)
			if score > 0.0 {
				matchingList = append(matchingList, matchingElem{aNodeInfo, score, bNodeInfo})
			}
		}
	}

	sort.Sort(sort.Reverse(matchingList))

	used := make(map[*nodeInfo]bool)
	matched := make(map[*nodeInfo]*nodeInfo)
	for _, elem := range matchingList {
		if _, ok := used[elem.match]; ok {
			continue
		}
		if _, ok := matched[elem.node]; ok {
			continue
		}
		used[elem.match] = true
		matched[elem.node] = elem.match
	}

	var result []matching
	j := 0
	for _, aNodeInfo := range aNodesInfo {
		bNodeInfo := matched[aNodeInfo]
		logrus.Debugln(callFunc, "matched:", reflect.TypeOf(aNodeInfo), aNodeInfo, "::",
			reflect.TypeOf(bNodeInfo), bNodeInfo)
		if bNodeInfo == nil {
			result = append(result, matching{
				prev:            aNodeInfo.node,
				positionChanged: false,
				orderChanged:    false,
			})
		} else {
			var orderChanged bool
			if j > bNodeInfo.pos && aNodeInfo.pos != bNodeInfo.pos {
				orderChanged = true
			}
			result = append(result, matching{
				prev:            aNodeInfo.node,
				next:            bNodeInfo.node,
				positionChanged: aNodeInfo.pos != bNodeInfo.pos,
				orderChanged:    orderChanged,
			})
			j = bNodeInfo.pos
		}
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

func matchSpecs(a, b []ast.Spec) []matching {
	return matchNodes(specToNodes(a), specToNodes(b), "matchIdents")
}

func specToNodes(l []ast.Spec) (nodes []ast.Node) {
	for _, e := range l {
		nodes = append(nodes, e)
	}
	return
}
