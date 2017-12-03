package diff

import "go/ast"

func diffDecl(a ast.Decl, b ast.Node, mode Mode) Coloring {
	b, ok := b.(ast.Decl)
	if !ok {
		return Coloring{NewColorChange(mode.ToColor(), a)}
	}
	switch t := a.(type) {
	case *ast.FuncDecl:
		return diffFuncDecl(t, b, mode)
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
