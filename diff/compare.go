package diff

import (
	"go/ast"
	"reflect"

	"math"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/gohist/util"
)

func compare(aNode, bNode ast.Node) (score float64) {
	defer func() { logrus.Debugln("compare:", "return:", score) }()
	if aNode == nil {
		if bNode == nil {
			return 1.0
		}
		return 0.0
	}
	switch a := aNode.(type) {
	case *ast.BadStmt:
		logrus.Debugln("comapare:", "*ast.BadStmt:", a, bNode)
		_, ok := bNode.(*ast.BadStmt)
		if ok {
			score += 1
		}
	case *ast.DeclStmt:
		logrus.Debugln("comapare:", "*ast.DeclStmt:", a, bNode)
		b, ok := bNode.(*ast.DeclStmt)
		if ok {
			score += compare(a.Decl, b.Decl)
		}
	case *ast.EmptyStmt:
		logrus.Debugln("comapare:", "*ast.EmptyStmt:", a, bNode)
		_, ok := bNode.(*ast.EmptyStmt)
		if ok {
			score += 1
		}
	case *ast.LabeledStmt:
		logrus.Debugln("comapare:", "*ast.LabeledStmt:", a, bNode)
		_, ok := bNode.(*ast.LabeledStmt)
		if ok {
			logrus.Errorln("compare:", "unimplemented:", reflect.TypeOf(a))
		}
	case *ast.ExprStmt:
		logrus.Debugln("comapare:", "*ast.ExprStmt:", a, bNode)
		b, ok := bNode.(*ast.ExprStmt)
		if ok {
			score += compare(a.X, b.X)
		}
	case *ast.SendStmt:
		logrus.Debugln("comapare:", "*ast.SendStmt:", a, bNode)
		b, ok := bNode.(*ast.SendStmt)
		if ok {
			score += compare(a.Chan, b.Chan) / 2
			score += compare(a.Value, b.Value) / 2
		}
	case *ast.IncDecStmt:
		logrus.Debugln("comapare:", "*ast.IncDecStmt:", a, bNode)
		b, ok := bNode.(*ast.IncDecStmt)
		if ok {
			score += compare(a.X, b.X) * (1 / math.Phi)
			if a.Tok == b.Tok {
				score += 1 - 1/math.Phi
			}
		}
	case *ast.AssignStmt:
		logrus.Debugln("comapare:", "*ast.AssignStmt:", a, bNode)
		b, ok := bNode.(*ast.AssignStmt)
		if ok {
			minLhs := util.IntMin(len(a.Lhs), len(b.Lhs))
			for i := 0; i < minLhs; i++ {
				score += compare(a.Lhs[i], b.Lhs[i])
			}
			minRhs := util.IntMin(len(a.Rhs), len(b.Rhs))
			for i := 0; i < minRhs; i++ {
				score += compare(a.Rhs[i], b.Rhs[i])
			}
			score = score / float64(minRhs+minLhs)
		}
	case *ast.GoStmt:
		logrus.Debugln("comapare:", "*ast.GoStmt:", a, bNode)
		b, ok := bNode.(*ast.GoStmt)
		if ok {
			score += compare(a.Call, b.Call)
		} else {
			//b, ok := bNode.(*ast.CallExpr)
			//if ok {
			//	score += compare(a.Call, b) * (1 - 1/math.Phi)
			//}
		}
	case *ast.DeferStmt:
		logrus.Debugln("comapare:", "*ast.DeferStmt:", a, bNode)
		b, ok := bNode.(*ast.DeferStmt)
		if ok {
			score += compare(a.Call, b.Call)
		}
	case *ast.ReturnStmt:
		logrus.Debugln("comapare:", "*ast.ReturnStmt:", a, bNode)
		b, ok := bNode.(*ast.ReturnStmt)
		if ok {
			if len(a.Results) == 0 && len(b.Results) == 0 {
				score = 1
			} else {
				max := util.IntMax(len(a.Results), len(b.Results))
				for _, match := range matchExprs(a.Results, b.Results) {
					if match.next != nil {
						score += compare(match.prev, match.next) / float64(max)
					}
				}
			}
		}
	case *ast.BranchStmt:
		logrus.Debugln("comapare:", "*ast.BranchStmt:", a, bNode)
		b, ok := bNode.(*ast.BranchStmt)
		if ok {
			if a.Tok == b.Tok {
				score += 1.0
			}
			if a.Label != nil {
				score = score / 2
				if b.Label != nil {
					score += compare(a.Label, b.Label) / 2
				}
			}
		}
	case *ast.BlockStmt:
		logrus.Debugln("comapare:", "*ast.BlockStmt:", a, bNode)
		b, ok := bNode.(*ast.BlockStmt)
		if ok {
			max := util.IntMax(len(a.List), len(b.List))
			for _, match := range matchStmts(a.List, b.List) {
				if match.next != nil {
					score += compare(match.prev, match.next) / float64(max)
				}
			}
		}
	case *ast.IfStmt:
		logrus.Debugln("comapare:", "*ast.IfStmt:", a, bNode)
		b, ok := bNode.(*ast.IfStmt)
		if ok {
			parts := 2.0
			if a.Init != nil {
				parts++
				score += compare(a.Init, b.Init)
			}
			score += compare(a.Cond, b.Cond)
			score += compare(a.Body, b.Body)
			if a.Else != nil {
				parts++
				score += compare(a.Else, b.Else)
			}
			score = score / parts
		}
	case *ast.SwitchStmt:
		logrus.Debugln("comapare:", "*ast.SwitchStmt:", a, bNode)
		b, ok := bNode.(*ast.SwitchStmt)
		if ok {
			score += compare(a.Init, b.Init) * (1 / math.Phi)
			score += compare(a.Body, b.Body) * (1 - 1/math.Phi)
		}
	case *ast.TypeSwitchStmt:
		logrus.Debugln("comapare:", "*ast.TypeSwitchStmt:", a, bNode)
		b, ok := bNode.(*ast.TypeSwitchStmt)
		if ok {
			score += compare(a.Assign, b.Assign) * (1 - 1/math.Phi)
			if a.Init != nil {
				score += compare(a.Init, b.Init) * (1 - 1/math.Phi)
				score = score / 2
			}
			score += compare(a.Body, b.Body) * (1 / math.Phi)
		}
	case *ast.SelectStmt:
		logrus.Debugln("comapare:", "*ast.SelectStmt:", a, bNode)
		b, ok := bNode.(*ast.SelectStmt)
		if ok {
			score += compare(a.Body, b.Body)
		}
	case *ast.ForStmt:
		logrus.Debugln("comapare:", "*ast.ForStmt:", a, bNode)
		b, ok := bNode.(*ast.ForStmt)
		if ok {
			children := 0
			if a.Init != nil {
				children++
				score += compare(a.Init, b.Init)
			}
			if a.Cond != nil {
				children++
				score += compare(a.Cond, b.Cond)
			}
			if a.Post != nil {
				children++
				score += compare(a.Post, b.Post)
			}
			if children > 0 {
				score = (score * (1 - 1/math.Phi)) / float64(children)
			}
			score += compare(a.Body, b.Body) / math.Phi
		}
	case *ast.RangeStmt:
		logrus.Debugln("comapare:", "*ast.RangeStmt:", a, bNode)
		b, ok := bNode.(*ast.RangeStmt)
		if ok {
			children := 1
			if a.Key != nil {
				children++
				score += compare(a.Key, b.Key)
			}
			if a.Value != nil {
				children++
				score += compare(a.Value, b.Value)
			}
			score += compare(a.X, b.X)
			score = (score * (1 - 1/math.Phi)) / float64(children)
			score += compare(a.Body, b.Body) / math.Phi
		}
	case *ast.Ident:
		logrus.Debugln("comapare:", "*ast.Ident:", a, bNode)
		b, ok := bNode.(*ast.Ident)
		if ok {
			if a.Name == b.Name {
				score += 1
			}
		}
	case *ast.CallExpr:
		logrus.Debugln("comapare:", "*ast.CallExpr:", a, bNode)
		b, ok := bNode.(*ast.CallExpr)
		if ok {
			total := util.IntMax(len(a.Args), len(b.Args))
			for _, match := range matchExprs(a.Args, b.Args) {
				if match.next != nil {
					score += compare(match.prev, match.next) / float64(total)
				}
			}
			score = score * (1 / math.Phi)
			score += compare(a.Fun, b.Fun) * (1 - (1 / math.Phi))
		}
	case *ast.StarExpr:
		logrus.Debugln("comapare:", "*ast.StarExpr:", a, bNode)
		b, ok := bNode.(*ast.StarExpr)
		if ok {
			score += compare(a.X, b.X)
		}
	case *ast.CaseClause:
		logrus.Debugln("comapare:", "*ast.CaseClause:", a, bNode)
		b, ok := bNode.(*ast.CaseClause)
		logrus.Debugln("comapare:", "(*ast.:", a, bNode)
		if ok {
			if len(a.List) == 0 && len(b.List) == 0 {
				score += 1
			} else {
				for _, match := range matchExprs(a.List, b.List) {
					if match.next != nil {
						score += compare(match.prev, match.next)
					}
				}
				score = score / float64(util.IntMax(len(a.List), len(b.List)))
			}

		}
	case *ast.SelectorExpr:
		logrus.Debugln("comapare:", "*ast.SelectorExpr:", a, bNode)
		b, ok := bNode.(*ast.SelectorExpr)
		if ok {
			score = compare(a.X, b.X) * (1 / math.Phi)
			if a.Sel.Name == b.Sel.Name {
				score += 1 - (1 / math.Phi)
			}
		}
	case *ast.BasicLit:
		logrus.Debugln("comapare:", "*ast.BasicLit:", a, bNode)
		b, ok := bNode.(*ast.BasicLit)
		if ok {
			if a.Kind == b.Kind {
				score += 0.5
				if a.Value == b.Value {
					score += 0.5
				}
			}
		}
	case *ast.TypeAssertExpr:
		logrus.Debugln("comapare:", "*ast.TypeAssertExpr:", a, bNode)
		b, ok := bNode.(*ast.TypeAssertExpr)
		if ok {
			score += compare(a.X, b.X)
			if a.Type != nil || b.Type != nil {
				score = (score + compare(a.Type, b.Type)) / 2
			}
		}
	case *ast.CompositeLit:
		logrus.Debugln("comapare:", "*ast.CompositeLit:", a, bNode)
		b, ok := bNode.(*ast.CompositeLit)
		if ok {
			score += compare(a.Type, b.Type)
		}
	case *ast.Field:
		logrus.Debugln("comapare:", "*ast.Field:", a, bNode)
		b, ok := bNode.(*ast.Field)
		if ok {
			score += compare(a.Type, b.Type)

		}
	case *ast.BinaryExpr:
		logrus.Debugln("comapare:", "*ast.BinaryExpr:", a, bNode)
		b, ok := bNode.(*ast.BinaryExpr)
		if ok {
			if a.Op == b.Op {
				score += 1 / 3
			}
			score += (compare(a.X, b.X) + compare(a.Y, b.Y)) / 3
		}
	case *ast.ArrayType:
		logrus.Debugln("comapare:", "*ast.ArrayType:", a, bNode)
		b, ok := bNode.(*ast.ArrayType)
		if ok {
			score += compare(a.Elt, b.Elt) * (1 / math.Phi)
			if a.Len != nil {
				score += compare(a.Len, b.Len) * (1 - (1 / math.Phi))
			} else {
				if b.Len == nil {
					score += 1 - 1/math.Phi
				}
			}
		}
	case *ast.FuncLit:
		logrus.Debugln("comapare:", "*ast.FuncLit:", a, bNode)
		b, ok := bNode.(*ast.FuncLit)
		if ok {
			score += compare(a.Type, b.Type) * (1 - 1/math.Phi)
			score += compare(a.Body, b.Body) * (1 / math.Phi)
		}
	case *ast.FuncType:
		logrus.Debugln("comapare:", "*ast.FuncType:", a, bNode)
		b, ok := bNode.(*ast.FuncType)
		if ok {
			score += compare(a.Params, b.Params) / 2
			score += compare(a.Results, b.Results) / 2
		}
	case *ast.FieldList:
		logrus.Debugln("comapare:", "*ast.FieldList:", a, bNode)
		b, ok := bNode.(*ast.FieldList)
		if ok {
			if a == nil {
				if b == nil {
					score += 1
				}
			} else if b != nil {
				max := util.IntMax(len(a.List), len(b.List))
				if max == 0 {
					score += 1
				}
				for _, match := range matchFields(a.List, b.List) {
					if match.next != nil {
						score += 1 / float64(max)
					}
				}
			}
		}
	case *ast.IndexExpr:
		logrus.Debugln("comapare:", "*ast.IndexExpr:", a, bNode)
		b, ok := bNode.(*ast.IndexExpr)
		if ok {
			score += compare(a.X, b.X) * 1 / math.Phi
			score += compare(a.Index, b.Index) * (1 - 1/math.Phi)
		}
	case *ast.MapType:
		logrus.Debugln("comapare:", "*ast.MapType:", a, bNode)
		b, ok := bNode.(*ast.MapType)
		if ok {
			score += compare(a.Key, b.Key) / 2
			score += compare(a.Value, b.Value) / 2
		}
	case *ast.GenDecl:
		logrus.Debugln("comapare:", "*ast.GenDecl:", a, bNode)
		b, ok := bNode.(*ast.GenDecl)
		if ok {
			max := util.IntMax(len(a.Specs), len(b.Specs))
			for _, match := range matchSpecs(a.Specs, b.Specs) {
				if match.next != nil {
					score += compare(match.prev, match.next) / float64(max)
				}
			}
		}
	case *ast.ValueSpec:
		logrus.Debugln("comapare:", "*ast.ValueSpec:", a, bNode)
		b, ok := bNode.(*ast.ValueSpec)
		if ok {
			max := util.IntMax(len(a.Names), len(b.Names))
			for _, match := range matchIdents(a.Names, b.Names) {
				if match.next != nil {
					score += compare(match.prev, match.next) / float64(max)
				}
			}
		}
	case *ast.ParenExpr:
		logrus.Debugln("comapare:", "*ast.ParenExpr:", a, bNode)
		b, ok := bNode.(*ast.ParenExpr)
		if ok {
			score = compare(a.X, b.X)
		}
	case *ast.SliceExpr:
		logrus.Debugln("comapare:", "*ast.SliceExpr:", a, bNode)
		b, ok := bNode.(*ast.SliceExpr)
		if ok {
			parts := 0
			if a.Low != nil {
				parts++
				score += compare(a.Low, b.Low)
			}
			if a.High != nil {
				parts++
				score += compare(a.High, b.High)
			}
			if a.Max != nil {
				parts++
				score += compare(a.Max, b.Max)
			}
			score = (score / float64(parts)) * (1 - 1/math.Phi)
			score += compare(a.X, b.X) * (1 / math.Phi)
		}
	case *ast.UnaryExpr:
		logrus.Debugln("comapare:", "*ast.UnaryExpr:", a, bNode)
		b, ok := bNode.(*ast.UnaryExpr)
		if ok {
			if a.Op == b.Op {
				score += 1 - 1/math.Phi
			}
			score += compare(a.X, b.X) * (1 / math.Phi)
		}
	case *ast.KeyValueExpr:
		logrus.Debugln("comapare:", "*ast.KeyValueExpr:", a, bNode)
		b, ok := bNode.(*ast.KeyValueExpr)
		if ok {
			score += (compare(a.Key, b.Key) + compare(a.Value, b.Value)) / 2
		}
	case *ast.InterfaceType:
		logrus.Debugln("comapare:", "*ast.InterfaceType:", a, bNode)
		b, ok := bNode.(*ast.InterfaceType)
		if ok {
			if a.Methods == nil {
				if b.Methods == nil {
					score = 1
				}
			} else {
				if b.Methods != nil {
					score += compare(a.Methods, b.Methods)
				}
			}
		}
	case *ast.ChanType:
		logrus.Debugln("comapare:", "*ast.ChanType:", a, bNode)
		b, ok := bNode.(*ast.ChanType)
		if ok {
			score += compare(a.Value, b.Value) * 1 / math.Phi
			if a.Dir == b.Dir {
				score += 1 - 1/math.Phi
			}
		}
	case *ast.CommClause:
		logrus.Debugln("comapare:", "*ast.CommClause:", a, bNode)
		b, ok := bNode.(*ast.CommClause)
		if ok {
			score += compare(a.Comm, b.Comm) * 1 / math.Phi
			max := float64(util.IntMax(len(a.Body), len(b.Body)))
			for _, match := range matchStmts(a.Body, b.Body) {
				if match.next != nil {
					score += compare(match.prev, match.prev) * (1 - 1/math.Phi) / max
				}
			}
		}
	case *ast.Ellipsis:
		logrus.Debugln("comapare:", "*ast.Ellipsis:", a, bNode)
		b, ok := bNode.(*ast.Ellipsis)
		if ok {
			score += compare(a.Elt, b.Elt)
		}
	default:
		logrus.Errorln("compare:", "unimplemented case: ", reflect.TypeOf(a))
	}
	return
}
