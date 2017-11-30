package diff

import (
	"go/ast"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func diffStmt(aStmt ast.Stmt, bNode ast.Node, mode Mode) Coloring {
	b, ok := bNode.(ast.Stmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), aStmt)}
	}
	switch a := aStmt.(type) {
	case *ast.BlockStmt:
		return diffBlockStmt(a, b, mode)
	case *ast.ForStmt:
		return diffForStmt(a, b, mode)
	case *ast.ExprStmt:
		return diffExprStmt(a, b, mode)
	case *ast.IfStmt:
		return diffIfStmt(a, b, mode)
	case *ast.AssignStmt:
		return diffAssignStmt(a, b, mode)
	case *ast.SwitchStmt:
		return diffSwitchStmt(a, b, mode)
	case *ast.CaseClause:
		return diffCaseClause(a, b, mode)
	default:
		logrus.Errorln("diffStmt:", "not implemented case", reflect.TypeOf(a))
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return nil
}

func diffBlockStmt(a *ast.BlockStmt, bNode ast.Node, mode Mode) Coloring {
	logrus.Debugln("diffBlockStmt:", a, bNode)
	b, ok := bNode.(*ast.BlockStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	var coloring Coloring
	for _, match := range matchStmts(a.List, b.List) {
		aStmt, bStmt := match.prev, match.next
		if bStmt == nil {
			logrus.Debugln("diffBlockStmt:", "unmatched:", aStmt, reflect.TypeOf(aStmt))
			coloring = append(coloring, NewColorChange(mode.ToColor(), aStmt))
		} else {
			coloring = append(coloring, diff(aStmt, bStmt, mode)...)
		}
	}
	return coloring
}

func diffForStmt(a *ast.ForStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffForStmt:", a, bNode)
	b, ok := bNode.(*ast.ForStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Init, b.Init, mode)...)
	coloring = append(coloring, diff(a.Cond, b.Cond, mode)...)
	coloring = append(coloring, diff(a.Post, b.Post, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)
	return coloring
}

func diffExprStmt(a *ast.ExprStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffExprStmt:", a, bNode)
	b, ok := bNode.(*ast.ExprStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diffExpr(a.X, b.X, mode)
}

func diffIfStmt(a *ast.IfStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffIfStmt:", a, bNode)
	b, ok := bNode.(*ast.IfStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Init, b.Init, mode)...)
	coloring = append(coloring, diff(a.Cond, b.Cond, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)
	return
}

func diffAssignStmt(a *ast.AssignStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffAssignStmt:", a, bNode)
	b, ok := bNode.(*ast.AssignStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	for _, match := range matchExprs(a.Lhs, b.Lhs) {
		aStmt, bStmt := match.prev, match.next
		if bStmt == nil {
			logrus.Debugln("diffAssignStmt:", "unmatched:", aStmt, reflect.TypeOf(aStmt))
			coloring = append(coloring, NewColorChange(mode.ToColor(), aStmt))
		} else {
			coloring = append(coloring, diff(aStmt, bStmt, mode)...)
		}
	}
	for _, match := range matchExprs(a.Rhs, b.Rhs) {
		aStmt, bStmt := match.prev, match.next
		if bStmt == nil {
			logrus.Debugln("diffAssignStmt:", "unmatched:", aStmt, reflect.TypeOf(aStmt))
			coloring = append(coloring, NewColorChange(mode.ToColor(), aStmt))
		} else {
			coloring = append(coloring, diff(aStmt, bStmt, mode)...)
		}
	}

	return
}

func diffSwitchStmt(a *ast.SwitchStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSwitchStmt:", a, bNode)
	b, ok := bNode.(*ast.SwitchStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Init, b.Init, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)

	return
}

func diffCaseClause(a *ast.CaseClause, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffCaseClause:", a, bNode)
	b, ok := bNode.(*ast.CaseClause)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	for _, match := range matchExprs(a.List, b.List) {
		if match.next == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}
	for _, match := range matchStmts(a.Body, b.Body) {
		if match.next == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}

	return
}