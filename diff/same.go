package diff

import (
	"go/ast"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func IsSame(aNode, bNode ast.Node) bool {
	if aNode == nil {
		return bNode == nil
	} else {
		if bNode == nil {
			return false
		}
	}
	switch a := aNode.(type) {
	case *ast.ArrayType:
		b, ok := bNode.(*ast.ArrayType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Elt, b.Elt) && IsSame(a.Len, b.Len)
	case *ast.AssignStmt:
		b, ok := bNode.(*ast.AssignStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Lhs) != len(b.Lhs) || len(a.Rhs) != len(b.Rhs) {
			return false
		}
		for i := 0; i < len(a.Lhs); i++ {
			if !IsSame(a.Lhs[i], b.Lhs[i]) {
				return false
			}
		}
		for i := 0; i < len(a.Rhs); i++ {
			if !IsSame(a.Rhs[i], b.Rhs[i]) {
				return false
			}
		}
		return true
	case *ast.BadDecl:
		b, ok := bNode.(*ast.BadDecl)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return true
	case *ast.BadExpr:
		b, ok := bNode.(*ast.BadExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return true
	case *ast.BadStmt:
		b, ok := bNode.(*ast.BadStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return true
	case *ast.BasicLit:
		b, ok := bNode.(*ast.BasicLit)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Kind == b.Kind && a.Value == b.Value
	case *ast.BinaryExpr:
		b, ok := bNode.(*ast.BinaryExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Op == b.Op && IsSame(a.X, b.X) && IsSame(a.Y, b.Y)
	case *ast.BlockStmt:
		b, ok := bNode.(*ast.BlockStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.List) != len(b.List) {
			return false
		}
		for i := 0; i < len(a.List); i++ {
			if !IsSame(a.List[i], b.List[i]) {
				return false
			}
		}
		return true
	case *ast.BranchStmt:
		b, ok := bNode.(*ast.BranchStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Tok == b.Tok && IsSame(a.Label, b.Label)
	case *ast.CallExpr:
		b, ok := bNode.(*ast.CallExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Args) != len(b.Args) {
			return false
		}
		for i := 0; i < len(a.Args); i++ {
			if !IsSame(a.Args[i], b.Args[i]) {
				return false
			}
		}
		return IsSame(a.Fun, b.Fun)
	case *ast.CaseClause:
		b, ok := bNode.(*ast.CaseClause)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.List) != len(b.List) || len(a.Body) != len(b.Body) {
			return false
		}
		for i := 0; i < len(a.List); i++ {
			if !IsSame(a.List[i], b.List[i]) {
				return false
			}
		}
		for i := 0; i < len(a.Body); i++ {
			if !IsSame(a.Body[i], b.Body[i]) {
				return false
			}
		}
		return true
	case *ast.ChanType:
		b, ok := bNode.(*ast.ChanType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Value, b.Value) && a.Dir == b.Dir
	case *ast.CommClause:
		b, ok := bNode.(*ast.CommClause)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Body) != len(b.Body) {
			return false
		}
		for i := 0; i < len(a.Body); i++ {
			if !IsSame(a.Body[i], b.Body[i]) {
				return false
			}
		}
		return IsSame(a.Comm, b.Comm)
	case *ast.Comment:
		_, ok := bNode.(*ast.Comment)
		if !ok {
			return false
		}
		return true
	case *ast.CommentGroup:
		_, ok := bNode.(*ast.CommentGroup)
		if !ok {
			return false
		}
		return true
	case *ast.CompositeLit:
		b, ok := bNode.(*ast.CompositeLit)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Elts) != len(b.Elts) {
			return false
		}
		for i := 0; i < len(a.Elts); i++ {
			if !IsSame(a.Elts[i], b.Elts[i]) {
				return false
			}
		}
		return IsSame(a.Type, b.Type)
	case *ast.DeclStmt:
		b, ok := bNode.(*ast.DeclStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Decl, b.Decl)
	case *ast.DeferStmt:
		b, ok := bNode.(*ast.DeferStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Call, b.Call)
	case *ast.Ellipsis:
		b, ok := bNode.(*ast.Ellipsis)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Elt, b.Elt)
	case *ast.EmptyStmt:
		b, ok := bNode.(*ast.EmptyStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return true
	case *ast.ExprStmt:
		b, ok := bNode.(*ast.ExprStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X)
	case *ast.Field:
		b, ok := bNode.(*ast.Field)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Names) != len(b.Names) {
			return false
		}
		for i := 0; i < len(a.Names); i++ {
			if !IsSame(a.Names[i], b.Names[i]) {
				return false
			}
		}
		return IsSame(a.Type, b.Type)
	case *ast.FieldList:
		b, ok := bNode.(*ast.FieldList)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.List) != len(b.List) {
			return false
		}
		for i := 0; i < len(a.List); i++ {
			if !IsSame(a.List[i], b.List[i]) {
				return false
			}
		}
		return true
	case *ast.ForStmt:
		b, ok := bNode.(*ast.ForStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Init, b.Init) && IsSame(a.Cond, b.Cond) && IsSame(a.Post, b.Post) && IsSame(a.Body, b.Body)
	case *ast.FuncDecl:
		b, ok := bNode.(*ast.FuncDecl)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Name, b.Name) && IsSame(a.Type, b.Type) && IsSame(a.Body, b.Body) // skip comments compare
	case *ast.FuncLit:
		b, ok := bNode.(*ast.FuncLit)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Type, b.Type) && IsSame(a.Body, b.Body)
	case *ast.FuncType:
		b, ok := bNode.(*ast.FuncType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Params, b.Params) && IsSame(a.Results, b.Results)
	case *ast.GenDecl:
		b, ok := bNode.(*ast.GenDecl)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Specs) != len(b.Specs) {
			return false
		}
		for i := 0; i < len(a.Specs); i++ {
			if !IsSame(a.Specs[i], b.Specs[i]) {
				return false
			}
		}
		return true
	case *ast.GoStmt:
		b, ok := bNode.(*ast.GoStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Call, b.Call)
	case *ast.Ident:
		b, ok := bNode.(*ast.Ident)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Name == b.Name // TODO: maybe improve that one
	case *ast.IfStmt:
		b, ok := bNode.(*ast.IfStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Init, b.Init) && IsSame(a.Cond, b.Cond) && IsSame(a.Body, b.Body) && IsSame(a.Else, b.Else)
	case *ast.ImportSpec:
		b, ok := bNode.(*ast.ImportSpec)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Name, b.Name)
	case *ast.IncDecStmt:
		b, ok := bNode.(*ast.IncDecStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Tok == b.Tok && IsSame(a.X, b.X)
	case *ast.IndexExpr:
		b, ok := bNode.(*ast.IndexExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X) && IsSame(a.Index, b.Index)
	case *ast.InterfaceType:
		b, ok := bNode.(*ast.InterfaceType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Methods, b.Methods)
	case *ast.KeyValueExpr:
		b, ok := bNode.(*ast.KeyValueExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Value, b.Value) && IsSame(a.Key, b.Key)
	case *ast.LabeledStmt:
		b, ok := bNode.(*ast.LabeledStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Label, b.Label) && IsSame(a.Stmt, b.Stmt)
	case *ast.MapType:
		b, ok := bNode.(*ast.MapType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Key, b.Key) && IsSame(a.Value, b.Value)
	case *ast.Package:
		b, ok := bNode.(*ast.Package)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Name == b.Name
	case *ast.ParenExpr:
		b, ok := bNode.(*ast.ParenExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X)
	case *ast.RangeStmt:
		b, ok := bNode.(*ast.RangeStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Key, b.Key) && IsSame(a.Value, b.Value) && IsSame(a.X, b.X)
	case *ast.ReturnStmt:
		b, ok := bNode.(*ast.ReturnStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Results) != len(b.Results) {
			return false
		}
		for i := 0; i < len(a.Results); i++ {
			if !IsSame(a.Results[i], b.Results[i]) {
				return false
			}
		}
		return true
	case *ast.SelectStmt:
		b, ok := bNode.(*ast.SelectStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Body, b.Body)
	case *ast.SelectorExpr:
		b, ok := bNode.(*ast.SelectorExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X) && IsSame(a.Sel, b.Sel)
	case *ast.SendStmt:
		b, ok := bNode.(*ast.SendStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Value, b.Value) && IsSame(a.Chan, b.Chan)
	case *ast.SliceExpr:
		b, ok := bNode.(*ast.SliceExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X) && IsSame(a.High, b.High) && IsSame(a.Low, b.Low) && IsSame(a.Max, b.Max)
	case *ast.StarExpr:
		b, ok := bNode.(*ast.StarExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X)
	case *ast.StructType:
		b, ok := bNode.(*ast.StructType)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Fields, b.Fields)
	case *ast.SwitchStmt:
		b, ok := bNode.(*ast.SwitchStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Init, b.Init) && IsSame(a.Tag, b.Tag) && IsSame(a.Body, b.Body)
	case *ast.TypeAssertExpr:
		b, ok := bNode.(*ast.TypeAssertExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.X, b.X) && IsSame(a.Type, b.Type)
	case *ast.TypeSpec:
		b, ok := bNode.(*ast.TypeSpec)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Type, b.Type) && IsSame(a.Name, b.Name)
	case *ast.TypeSwitchStmt:
		b, ok := bNode.(*ast.TypeSwitchStmt)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return IsSame(a.Assign, b.Assign) && IsSame(a.Init, b.Init) && IsSame(a.Body, b.Body)
	case *ast.UnaryExpr:
		b, ok := bNode.(*ast.UnaryExpr)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		return a.Op == b.Op && IsSame(a.X, b.X)
	case *ast.ValueSpec:
		b, ok := bNode.(*ast.ValueSpec)
		if !ok {
			return false
		}
		if a == nil {
			return b == nil
		} else {
			if b == nil {
				return false
			}
		}
		if len(a.Names) != len(b.Names) || len(a.Values) != len(b.Values) {
			return false
		}
		for i := 0; i < len(a.Names); i++ {
			if !IsSame(a.Names[i], b.Names[i]) {
				return false
			}
		}
		for i := 0; i < len(a.Values); i++ {
			if !IsSame(a.Values[i], b.Values[i]) {
				return false
			}
		}
		return IsSame(a.Type, b.Type)
	default:
		logrus.Errorln("unimplemented case:", reflect.TypeOf(a), reflect.TypeOf(bNode))
		return false
	}
}
