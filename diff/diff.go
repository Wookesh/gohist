package diff

import (
	"go/ast"
	"go/token"
	"reflect"

	"github.com/Sirupsen/logrus"
)

type context struct {
	a     nodeContext
	b     nodeContext
	gobal vars
}

type nodeContext struct {
	vars vars
}

type vars map[string][]token.Pos

func Diff(a, b ast.Node, mode Mode) Coloring {
	logrus.Debugln("Diff:", mode, "\n")
	if mode == ModeNew && a == nil {
		return Coloring{NewColorChange(mode.ToColor(), b)}
	}
	if mode == ModeOld && b == nil {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a, b, mode)
}

func diff(aNode, b ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diff:", aNode, b)
	if aNode == nil && b == nil {
		return
	}
	switch a := aNode.(type) {
	case ast.Decl:
		coloring = diffDecl(a, b, mode)
	case ast.Stmt:
		coloring = diffStmt(a, b, mode)
	case ast.Expr:
		coloring = diffExpr(a, b, mode)
		//case ast.Decl:
		//	diffDecl(t, b, mode)
	// non interface nodes:
	case *ast.FieldList:
		coloring = diffFieldList(a, b, mode)
	case *ast.Field:
		coloring = diffField(a, b, mode)
	default:
		logrus.Errorln("diff:", "not implemented case", reflect.TypeOf(a))
		coloring = Coloring{NewColorChange(mode.ToColor(), a)}
	}

	return
}

func diffFieldList(a *ast.FieldList, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.FieldList)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a == nil {
		return
	}
	for _, match := range matchFields(a.List, b.List) {
		if match.next == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}
	return
}

func diffField(a *ast.Field, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.Field)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}

	for _, match := range matchIdents(a.Names, b.Names) {
		if match.next == nil {
			coloring = append(coloring, NewColorChange(mode.ToColor(), match.prev))
		} else {
			coloring = append(coloring, diff(match.prev, match.next, mode)...)
		}
	}
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	return
}
