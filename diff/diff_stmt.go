package diff

import (
	"go/ast"
	"reflect"

	"github.com/sirupsen/logrus"
)

func diffStmt(aStmt ast.Stmt, bNode ast.Node, mode Mode) Coloring {
	b, ok := bNode.(ast.Stmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), aStmt)}
	}
	switch a := aStmt.(type) {
	case *ast.AssignStmt:
		return diffAssignStmt(a, b, mode)
	case *ast.BlockStmt:
		return diffBlockStmt(a, b, mode)
	case *ast.BranchStmt:
		return diffBranchStmt(a, b, mode)
	case *ast.CaseClause:
		return diffCaseClause(a, b, mode)
	case *ast.CommClause:
		return diffCommClause(a, b, mode)
	case *ast.DeclStmt:
		return diffDeclStmt(a, b, mode)
	case *ast.DeferStmt:
		return diffDeferStmt(a, b, mode)
	case *ast.ExprStmt:
		return diffExprStmt(a, b, mode)
	case *ast.ForStmt:
		return diffForStmt(a, b, mode)
	case *ast.GoStmt:
		return diffGoStmt(a, b, mode)
	case *ast.IfStmt:
		return diffIfStmt(a, b, mode)
	case *ast.IncDecStmt:
		return diffIncDecStmt(a, b, mode)
	case *ast.RangeStmt:
		return diffRangeStmt(a, b, mode)
	case *ast.ReturnStmt:
		return diffReturnStmt(a, b, mode)
	case *ast.SelectStmt:
		return diffSelectStmt(a, b, mode)
	case *ast.SendStmt:
		return diffSendStmt(a, b, mode)
	case *ast.SwitchStmt:
		return diffSwitchStmt(a, b, mode)
	case *ast.TypeSwitchStmt:
		return diffTypeSwitchStmt(a, b, mode)
	default:
		logrus.Errorln("diffStmt:", "not implemented case", reflect.TypeOf(a))
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return nil
}

func diffBlockStmt(a *ast.BlockStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffBlockStmt:", a, bNode)
	b, ok := bNode.(*ast.BlockStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, colorMatches(matchStmts(a.List, b.List), mode, "diffBlockStmt")...)
	return
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
	if a.Else != nil {
		coloring = append(coloring, diff(a.Else, b.Else, mode)...)
	}
	return
}

func diffAssignStmt(a *ast.AssignStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffAssignStmt:", a, bNode)
	b, ok := bNode.(*ast.AssignStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, colorMatches(matchExprs(a.Lhs, b.Lhs), mode, "diffAssignStmt")...)
	coloring = append(coloring, colorMatches(matchExprs(a.Rhs, b.Rhs), mode, "diffAssignStmt")...)

	return
}

func diffSwitchStmt(a *ast.SwitchStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSwitchStmt:", a, bNode)
	b, ok := bNode.(*ast.SwitchStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Init, b.Init, mode)...)
	coloring = append(coloring, diff(a.Tag, b.Tag, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)

	return
}

func diffCaseClause(a *ast.CaseClause, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffCaseClause:", a, bNode)
	b, ok := bNode.(*ast.CaseClause)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, colorMatches(matchExprs(a.List, b.List), mode, "diffCaseClause")...)
	coloring = append(coloring, colorMatches(matchStmts(a.Body, b.Body), mode, "diffCaseClause")...)
	return
}

func diffDeclStmt(a *ast.DeclStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffDeclStmt:", a, bNode)
	b, ok := bNode.(*ast.DeclStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diffDecl(a.Decl, b.Decl, mode)
}

func diffTypeSwitchStmt(a *ast.TypeSwitchStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffTypeSwitchStmt:", a, bNode)
	b, ok := bNode.(*ast.TypeSwitchStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	coloring = append(coloring, diff(a.Assign, b.Assign, mode)...)
	coloring = append(coloring, diff(a.Init, b.Init, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)
	return
}

func diffReturnStmt(a *ast.ReturnStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffReturnStmt:", a, bNode)
	b, ok := bNode.(*ast.ReturnStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, colorMatches(matchExprs(a.Results, b.Results), mode, "diffReturnStmt")...)
	return
}

func diffRangeStmt(a *ast.RangeStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffRangeStmt:", a, bNode)
	b, ok := bNode.(*ast.RangeStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a.Key != nil {
		coloring = append(coloring, diff(a.Key, b.Key, mode)...)
	}
	if a.Value != nil {
		coloring = append(coloring, diff(a.Value, b.Value, mode)...)
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)
	return
}

func diffIncDecStmt(a *ast.IncDecStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffIncDecStmt:", a, bNode)
	b, ok := bNode.(*ast.IncDecStmt)
	if !ok || a.Tok != b.Tok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a.X, b.X, mode)
}

func diffBranchStmt(a *ast.BranchStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffBranchStmt:", a, bNode)
	b, ok := bNode.(*ast.BranchStmt)
	if !ok || a.Tok != b.Tok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a.Label != nil {
		if b.Label != nil {
			return diff(a.Label, b.Label, mode)
		} else {
			return Coloring{NewColorChange(mode.ToColor(), a)}
		}
	}
	return
}

func diffGoStmt(a *ast.GoStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffGoStmt:", a, bNode)
	b, ok := bNode.(*ast.GoStmt)
	if !ok {
		//b, ok := bNode.(*ast.CallExpr)
		//if !ok {
		//	return Coloring{NewColorChange(mode.ToColor(), a)}
		//}
		//coloring = append(coloring, ColorChange{Color: mode.ToColor(), Pos: a.Go, End: a.Pos()})
		//coloring = append(coloring, diff(a.Call, b, mode)...)
		//return
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = diff(a.Call, b.Call, mode)
	return
}

func diffDeferStmt(a *ast.DeferStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffDeferStmt:", a, bNode)
	b, ok := bNode.(*ast.DeferStmt)
	if !ok {
		//b, ok := bNode.(*ast.CallExpr)
		//if !ok {
		//	return Coloring{NewColorChange(mode.ToColor(), a)}
		//}
		//coloring = append(coloring, ColorChange{Color: mode.ToColor(), Pos: a.Defer, End: a.Pos()})
		//coloring = append(coloring, diff(a.Call, b, mode)...)
		//return
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = diff(a.Call, b.Call, mode)
	return
}

func diffSelectStmt(a *ast.SelectStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSelectStmt:", a, bNode)
	b, ok := bNode.(*ast.SelectStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = diff(a.Body, b.Body, mode)
	return
}

func diffCommClause(a *ast.CommClause, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffCommClause:", a, bNode)
	b, ok := bNode.(*ast.CommClause)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Comm, b.Comm, mode)...)
	coloring = append(coloring, colorMatches(matchStmts(a.Body, b.Body), mode, "CommClause")...)
	return
}

func diffSendStmt(a *ast.SendStmt, bNode ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSendStmt:", a, bNode)
	b, ok := bNode.(*ast.SendStmt)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Chan, b.Chan, mode)...)
	coloring = append(coloring, diff(a.Value, b.Value, mode)...)
	return
}
