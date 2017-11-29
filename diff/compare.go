package diff

import (
	"fmt"
	"go/ast"
	"reflect"

	"github.com/Sirupsen/logrus"
)

func compare(aNode, bNode ast.Node) (score float64) {
	switch a := aNode.(type) {
	case *ast.BadStmt:
		_, ok := bNode.(*ast.BadStmt)
		if ok {
			score += 1
		}
	case *ast.DeclStmt:
		_, ok := bNode.(*ast.DeclStmt)
		if ok {
			score += 1 //compare(a.Decl, b.Decl)
		}
	case *ast.EmptyStmt:
		_, ok := bNode.(*ast.EmptyStmt)
		if ok {
			score += 1
		}
	case *ast.LabeledStmt:
		_, ok := bNode.(*ast.LabeledStmt)
		if ok {
			score += 1
		}
	case *ast.ExprStmt:
		b, ok := bNode.(*ast.ExprStmt)
		if ok {
			score += compareExpr(a.X, b.X)
		}
	case *ast.SendStmt:
		_, ok := bNode.(*ast.SendStmt)
		if ok {
			score += 1
		}
	case *ast.IncDecStmt:
		_, ok := bNode.(*ast.IncDecStmt)
		if ok {
			score += 1
		}
	case *ast.AssignStmt:
		_, ok := bNode.(*ast.AssignStmt)
		if ok {
			score += 1
		}
	case *ast.GoStmt:
		_, ok := bNode.(*ast.GoStmt)
		if ok {
			score += 1
		}
	case *ast.DeferStmt:
		_, ok := bNode.(*ast.DeferStmt)
		if ok {
			score += 1
		}
	case *ast.ReturnStmt:
		_, ok := bNode.(*ast.ReturnStmt)
		if ok {
			score += 1
		}
	case *ast.BranchStmt:
		_, ok := bNode.(*ast.BranchStmt)
		if ok {
			score += 1
		}
	case *ast.BlockStmt:
		_, ok := bNode.(*ast.BlockStmt)
		if ok {
			score += 1
		}
	case *ast.IfStmt:
		_, ok := bNode.(*ast.IfStmt)
		if ok {
			score += 1
		}
	case *ast.SwitchStmt:
		_, ok := bNode.(*ast.SwitchStmt)
		if ok {
			score += 1
		}
	case *ast.TypeSwitchStmt:
		_, ok := bNode.(*ast.TypeSwitchStmt)
		if ok {
			score += 1
		}
	case *ast.SelectStmt:
		_, ok := bNode.(*ast.SelectStmt)
		if ok {
			score += 1
		}
	case *ast.ForStmt:
		b, ok := bNode.(*ast.ForStmt)
		if ok {
			score += 1
			fmt.Println("for:", a, b)
		}
	case *ast.RangeStmt:
		b, ok := bNode.(*ast.RangeStmt)
		if ok {
			score += 1
			fmt.Println("range:", a, b)
		}
	default:
		logrus.Errorln("compare:", "unimplemented case: ", reflect.TypeOf(a))
	}
	return
}

func compareExpr(aExpr, bExpr ast.Expr) (score float64) {
	logrus.Debugln("compareExpr:", aExpr, reflect.TypeOf(aExpr), bExpr, reflect.TypeOf(bExpr))
	switch a := aExpr.(type) {
	case *ast.CallExpr:
		b, ok := bExpr.(*ast.CallExpr)
		if ok {
			score = compareExpr(a.Fun, b.Fun)
		}
	case *ast.SelectorExpr:
		b, ok := bExpr.(*ast.SelectorExpr)
		if ok {
			score = compareExpr(a.X, b.X)
			if a.Sel.Name == b.Sel.Name {
				score += 1
			}
			score = score / 2
		}
	case *ast.Ident:
		b, ok := bExpr.(*ast.Ident)
		if ok {
			if a.Name == b.Name {
				score += 1
			}
		}
	default:
		fmt.Println("compareExpr:", "unimplemented case: ", reflect.TypeOf(a))
	}
	return
}
