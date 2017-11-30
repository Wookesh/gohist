package diff

import "go/ast"

func diffDecl(a ast.Decl, b ast.Node, mode Mode) Coloring {
	_, ok := b.(ast.Decl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	switch t := a.(type) {
	case *ast.FuncDecl:
		return diffFuncDecl(t, b, mode)
	}
	return nil
}

func diffFuncDecl(a *ast.FuncDecl, bNode ast.Node, mode Mode) Coloring {
	b, ok := bNode.(*ast.FuncDecl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	return diff(a.Body, b.Body, mode)
}