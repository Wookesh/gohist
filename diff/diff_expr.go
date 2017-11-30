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
	case *ast.BasicLit:
		return diffBasicLit(a, bExpr, mode)
	case *ast.TypeAssertExpr:
		return diffTypeAssertExpr(a, bExpr, mode)
	case *ast.CompositeLit:
		return diffCompositeLit(a, bExpr, mode)
	case *ast.FuncType:
		return diffFuncType(a, bExpr, mode)
	case *ast.UnaryExpr:
		return diffUnaryExpr(a, bExpr, mode)
	case *ast.ArrayType:
		return diffArrayType(a, bExpr, mode)
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
	coloring = append(coloring, diffExpr(a.X, b.X, mode)...)
	if a.Sel.Name != b.Sel.Name {
		coloring = append(coloring, NewColorChange(mode.ToColor(), a.Sel))
	}
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

func diffBasicLit(a *ast.BasicLit, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffBasicLit:", a, bExpr)
	b, ok := bExpr.(*ast.BasicLit)
	if !ok || a.Kind != b.Kind || a.Value != b.Value {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return
}

func diffTypeAssertExpr(a *ast.TypeAssertExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffTypeAssertExpr:", a, bExpr)
	b, ok := bExpr.(*ast.TypeAssertExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	return
}

func diffCompositeLit(a *ast.CompositeLit, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffCompositeLit:", a, bExpr)
	b, ok := bExpr.(*ast.CompositeLit)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	for _, match := range matchExprs(a.Elts, b.Elts) {
		if match.next == nil {
			logrus.Debugln("diffCompositeLit:", "unmatched:", match.prev, reflect.TypeOf(match.prev))
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}
	return
}

func diffFuncType(a *ast.FuncType, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffFuncType:", a, bExpr)
	b, ok := bExpr.(*ast.FuncType)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Params, b.Params, mode)...)
	coloring = append(coloring, diff(a.Results, b.Results, mode)...)
	return
}

func diffUnaryExpr(a *ast.UnaryExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffUnaryExpr:", a, bExpr)
	b, ok := bExpr.(*ast.UnaryExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	return
}

func diffArrayType(a *ast.ArrayType, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffArrayType:", a, bExpr)
	b, ok := bExpr.(*ast.ArrayType)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Len, b.Len, mode)...)
	coloring = append(coloring, diff(a.Elt, b.Elt, mode)...)
	return
}
