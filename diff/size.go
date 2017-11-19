package diff

import "go/ast"

func getSize(a ast.Node) (size float64) {
	size = 1
	switch t := a.(type) {
	case *ast.BadExpr, *ast.Ident, *ast.Ellipsis, *ast.BasicLit, *ast.FuncLit, *ast.CompositeLit:
		break
	case *ast.ParenExpr:
		size += getSize(t.X)
	case *ast.SelectorExpr:
		size += getSize(t.X)
	case *ast.IndexExpr:
		size += getSize(t.X)
	case *ast.SliceExpr:
		size += getSize(t.X)
	case *ast.TypeAssertExpr:
		size += getSize(t.X) + getSize(t.Type)
	case *ast.CallExpr:
		for _, arg := range t.Args {
			size += getSize(arg)
		}
		size += getSize(t.Fun)
	case *ast.StarExpr:
		size += getSize(t.X)
	case *ast.UnaryExpr:
		size += getSize(t.X)
	case *ast.BinaryExpr:
		size += getSize(t.X) + getSize(t.Y)
	case *ast.KeyValueExpr:
		size += getSize(t.Key) + getSize(t.Value)
	case *ast.ArrayType:
		size += getSize(t.Elt)
	case *ast.StructType:
		panic("getSize, case: *ast.StructType")
	case *ast.FuncType:
		panic("getSize, case: *ast.FuncType")
	case *ast.InterfaceType:
		panic("getSize, case: *ast.InterfaceType")
	case *ast.MapType:
		panic("getSize, case: *ast.MapType")
	case *ast.ChanType:
		panic("getSize, case: *ast.ChanType")
	case *ast.BadStmt:
		panic("getSize, case: *ast.BadStmt")
	case *ast.DeclStmt:
		size += getSize(t.Decl)
	case *ast.EmptyStmt:
		panic("getSize, case: *ast.EmptyStmt")
	case *ast.LabeledStmt:
		panic("getSize, case: *ast.LabeledStmt")
	case *ast.ExprStmt:
		size += getSize(t.X)
	case *ast.SendStmt:
		size += getSize(t.Chan) + getSize(t.Value)
	case *ast.IncDecStmt:
		panic("getSize, case: *ast.IncDecStmt")
	case *ast.AssignStmt:
		for _, expr := range t.Lhs {
			size += getSize(expr)
		}
		for _, expr := range t.Rhs {
			size += getSize(expr)
		}
	case *ast.GoStmt:
		panic("getSize, case: *ast.GoStmt")
	case *ast.DeferStmt:
		panic("getSize, case: *ast.DeferStmt")
	case *ast.ReturnStmt:
		panic("getSize, case: *ast.ReturnStmt")
	case *ast.BranchStmt:
		panic("getSize, case: *ast.BranchStmt")
	case *ast.BlockStmt:
		for _, stmt := range t.List {
			size += getSize(stmt)
		}
	case *ast.IfStmt:
		size += getSize(t.Body) + getSize(t.Cond) + getSize(t.Else) + getSize(t.Init)
	case *ast.CaseClause:
		panic("getSize, case: *ast.CaseClause")
	case *ast.SwitchStmt:
		panic("getSize, case: *ast.SwitchStmt")
	case *ast.TypeSwitchStmt:
		panic("getSize, case: *ast.TypeSwitchStmt")
	case *ast.CommClause:
		panic("getSize, case: *ast.CommClause")
	case *ast.SelectStmt:
		panic("getSize, case: *ast.SelectStmt")
	case *ast.ForStmt:
		panic("getSize, case: *ast.ForStmt")
	case *ast.RangeStmt:
		size += getSize(t.Body) + getSize(t.Value) + getSize(t.X) + getSize(t.Key)
	case *ast.ImportSpec:
		panic("getSize, case: *ast.ImportSpec")
	case *ast.ValueSpec:
		panic("getSize, case: *ast.ValueSpec")
	case *ast.TypeSpec:
		panic("getSize, case: *ast.TypeSpec")
	case *ast.BadDecl:
		panic("getSize, case: *ast.BadDecl")
	case *ast.GenDecl:
		panic("getSize, case: *ast.GenDecl")
	case *ast.FuncDecl:
		panic("getSize, case: *ast.FuncDecl")

	}

	return
}