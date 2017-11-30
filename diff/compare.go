package diff

import (
	"go/ast"
	"reflect"

	"math"

	"github.com/Sirupsen/logrus"
	"github.com/wookesh/gohist/util"
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
			score += compare(a.X, b.X)
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
		_, ok := bNode.(*ast.ForStmt)
		if ok {
			score += 1
		}
	case *ast.RangeStmt:
		_, ok := bNode.(*ast.RangeStmt)
		if ok {
			score += 1
		}
	case *ast.Ident:
		b, ok := bNode.(*ast.Ident)
		if ok {
			if a.Name == b.Name {
				score += 1
			}
		}
	case *ast.CallExpr:
		b, ok := bNode.(*ast.CallExpr)
		if ok {
			score += compare(a.Fun, b.Fun)
		}
	case *ast.StarExpr:
		b, ok := bNode.(*ast.StarExpr)
		if ok {
			score += compare(a.X, b.X)
		}
	case *ast.CaseClause:
		b, ok := bNode.(*ast.CaseClause)
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
		b, ok := bNode.(*ast.SelectorExpr)
		if ok {
			score = compare(a.X, b.X) * (1 / math.Phi)
			if a.Sel.Name == b.Sel.Name {
				score += 1 * (1 - (1 / math.Phi))
			}
		}
	case *ast.BasicLit:
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
		b, ok := bNode.(*ast.TypeAssertExpr)
		if ok {
			score += compare(a.X, b.X)
			if a.Type != nil || b.Type != nil {
				score = (score + compare(a.Type, b.Type)) / 2
			}
		}
	case *ast.CompositeLit:
		b, ok := bNode.(*ast.CompositeLit)
		if ok {
			score += compare(a.Type, b.Type)
		}
	case *ast.Field:
		b, ok := bNode.(*ast.Field)
		if ok {
			score += compare(a.Type, b.Type)

		}
	case *ast.BinaryExpr:
		b, ok := bNode.(*ast.BinaryExpr)
		if ok {
			if a.Op == b.Op {
				score += 1 / 3
			}
			score += (compare(a.X, b.X) + compare(a.Y, b.Y)) / 3
		}
	default:
		logrus.Errorln("compare:", "unimplemented case: ", reflect.TypeOf(a))
	}
	return
}
