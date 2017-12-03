package diff

import (
	"go/ast"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func diffDecl(aDecl ast.Decl, bDecl ast.Node, mode Mode) Coloring {
	b, ok := bDecl.(ast.Decl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), aDecl)}
	}
	switch a := aDecl.(type) {
	case *ast.FuncDecl:
		return diffFuncDecl(a, b, mode)
	case *ast.GenDecl:
		return diffGenDecl(a, b, mode)
	default:
		logrus.Errorln("diffDecl:", "unimplemented case:", reflect.TypeOf(a))
	}
	return nil
}

func diffFuncDecl(a *ast.FuncDecl, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.FuncDecl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a == nil {
		return nil
	} else {
		if b == nil {
			return Coloring{NewColorChange(mode.ToColor(), a)}
		}
	}
	coloring = append(coloring, diff(a.Type, b.Type, mode)...)
	coloring = append(coloring, diff(a.Name, b.Name, mode)...)
	coloring = append(coloring, diff(a.Body, b.Body, mode)...)

	return
}

func diffGenDecl(a *ast.GenDecl, bNode ast.Node, mode Mode) (coloring Coloring) {
	b, ok := bNode.(*ast.GenDecl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	if a == nil {
		return nil
	} else {
		if b == nil {
			return Coloring{NewColorChange(mode.ToColor(), a)}
		}
	}
	if a.Tok != b.Tok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	coloring = append(coloring, colorMatches(matchSpecs(a.Specs, b.Specs), mode, "diffGenDecl")...)
	return
}
