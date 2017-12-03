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
	case *ast.FuncLit:
		return diffFuncLit(a, bExpr, mode)
	case *ast.IndexExpr:
		return diffIndexExpr(a, bExpr, mode)
	case *ast.MapType:
		return diffMapType(a, bExpr, mode)
	case *ast.ParenExpr:
		return diffParenExpr(a, bExpr, mode)
	case *ast.SliceExpr:
		return diffSliceExpr(a, bExpr, mode)
	case *ast.KeyValueExpr:
		return diffKeyValueExpr(a, bExpr, mode)
	case *ast.InterfaceType:
		return diffInterfaceType(a, bExpr, mode)
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
	//if len(coloring) > 0 {
	//	return Coloring{NewColorChange(mode.ToColor(), a)}
	//}
	coloring = append(coloring, colorMatches(matchExprs(a.Args, b.Args), mode, "diffCallExpr")...)
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
	coloring = append(coloring, colorMatches(matchExprs(a.Elts, b.Elts), mode, "diffCompositeLit")...)
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

func diffFuncLit(a *ast.FuncLit, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffFuncLit:", a, bExpr)
	b, ok := bExpr.(*ast.FuncLit)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)
	return
}

func diffIndexExpr(a *ast.IndexExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffIndexExpr:", a, bExpr)
	b, ok := bExpr.(*ast.IndexExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	coloring = append(coloring, diff(a.Index, b.Index, mode)...)
	return
}

func diffMapType(a *ast.MapType, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffMapType:", a, bExpr)
	b, ok := bExpr.(*ast.MapType)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Key, b.Key, mode)...)
	coloring = append(coloring, diff(a.Value, b.Value, mode)...)
	return
}

func diffParenExpr(a *ast.ParenExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffParenExpr:", a, bExpr)
	b, ok := bExpr.(*ast.ParenExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	return
}

func diffSliceExpr(a *ast.SliceExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffSliceExpr:", a, bExpr)
	b, ok := bExpr.(*ast.SliceExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.X, b.X, mode)...)
	if a.High != nil {
		coloring = append(coloring, diff(a.High, b.High, mode)...)
	}
	if a.Low != nil {
		coloring = append(coloring, diff(a.Low, b.Low, mode)...)
	}
	if a.Max != nil {
		coloring = append(coloring, diff(a.Max, b.Max, mode)...)
	}
	return
}

func diffKeyValueExpr(a *ast.KeyValueExpr, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffKeyValueExpr:", a, bExpr)
	b, ok := bExpr.(*ast.KeyValueExpr)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, diff(a.Key, b.Key, mode)...)
	coloring = append(coloring, diff(a.Value, b.Value, mode)...)

	return
}

func diffInterfaceType(a *ast.InterfaceType, bExpr ast.Expr, mode Mode) (coloring Coloring) {
	logrus.Debugln("diffInterfaceType:", a, bExpr)
	b, ok := bExpr.(*ast.InterfaceType)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = diff(a.Methods, b.Methods, mode)
	return
}
