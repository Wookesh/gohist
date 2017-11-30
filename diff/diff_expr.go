package diff

import (
	"go/ast"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func diffExpr(aExpr ast.Expr, bNode ast.Node, mode Mode) Coloring {
	logrus.Debugln("diffExpr:", aExpr, bNode)
	bExpr, ok := bNode.(ast.Expr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), aExpr)}
	}
	switch a := aExpr.(type) {
	case *ast.CallExpr:
		return diffCallExpr(a, bExpr, mode)
	case *ast.SelectorExpr:
		return diffSelectorExpr(a, bExpr, mode)
	case *ast.Ident:
		return diffIdent(a, bExpr, mode)
	case *ast.BinaryExpr:
		return diffBinaryExpr(a, bExpr, mode)
	case *ast.StarExpr:
		return diffStarExpr(a, bExpr, mode)
	default:
		logrus.Errorln("diffExpr:", "unimplemented case:", reflect.TypeOf(a))
		return Coloring{NewColorChange(mode.ToColor(), aExpr)}
	}
}

func diffCallExpr(a *ast.CallExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	b, ok := bExpr.(*ast.CallExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = diffExpr(a.Fun, b.Fun, mode)
	if len(coloring) > 0 {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	for _, match := range matchExprs(a.Args, b.Args) {
		aExpr, bExpr := match.prev, match.next
		if match.next == nil {
			logrus.Debugln("diffCallExpr:", "unmatched:", aExpr, reflect.TypeOf(aExpr))
			coloring = append(coloring, NewColorChange(mode.ToColor(), aExpr))
		} else {
			coloring = append(coloring, diff(aExpr, bExpr, mode)...)
		}
	}
	return coloring
}

func diffSelectorExpr(a *ast.SelectorExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSelectorExpr:", a, bExpr)
	b, ok := bExpr.(*ast.SelectorExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a.Sel.Name != b.Sel.Name {
		coloring = append(coloring, NewColorChange(mode.ToColor(), a.Sel))
	}
	coloring = append(coloring, diffExpr(a.X, b.X, mode)...)
	return
}

func diffIdent(a *ast.Ident, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffIdent:", a, bExpr)
	b, ok := bExpr.(*ast.Ident)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a.Name != b.Name {
		coloring = append(coloring, NewColorChange(mode.ToColor(), a))
	}
	return
}

func diffBinaryExpr(a *ast.BinaryExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffBinaryExpr:", a, bExpr)
	b, ok := bExpr.(*ast.BinaryExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a.Op != b.Op {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	coloring = append(coloring, diff(a.Y, b.Y, mode)...)
	return
}

func diffStarExpr(a *ast.StarExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffStarExpr:", a, bExpr)
	b, ok := bExpr.(*ast.StarExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a.X, b.X, mode)
}
