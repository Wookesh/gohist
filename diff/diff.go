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

func diff(a, b ast.Node, mode Mode) (coloring Coloring) {
	logrus.Debugln("diff:", a, b)
	if a == nil && b == nil {
		return
	}
	switch t := a.(type) {
	case ast.Decl:
		coloring = diffDecl(t, b, mode)
	case ast.Stmt:
		coloring = diffStmt(t, b, mode)
	case ast.Expr:
		coloring = diffExpr(t, b, mode)
		//case ast.Decl:
		//	diffDecl(t, b, mode)
	default:
		logrus.Errorln("diff:", "not implemented case", reflect.TypeOf(a))
		coloring = Coloring{NewColorChange(mode.ToColor(), a)}
	}

	return
}
