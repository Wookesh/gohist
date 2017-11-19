package diff

import (
	"go/ast"
	"log"
)

func compare(a, b ast.Node) (score float64) {
	switch t := a.(type) {
	case *ast.BadStmt:
		_, ok := b.(*ast.BadStmt)
		if ok {
			score += 1
		}
	case *ast.DeclStmt:
		bT, ok := b.(*ast.DeclStmt)
		if !ok {
			return
		}
		log.Print(t)
		score += compare(t.Decl, bT.Decl)
	case *ast.EmptyStmt:
		_, ok := b.(*ast.EmptyStmt)
		if ok {
			score += 1
		}
	case *ast.LabeledStmt:
		_, ok := b.(*ast.LabeledStmt)
		if ok {
			score += 1
		}
	case *ast.ExprStmt:
		_, ok := b.(*ast.ExprStmt)
		if ok {
			score += 1
		}
	case *ast.SendStmt:
		_, ok := b.(*ast.SendStmt)
		if ok {
			score += 1
		}
	case *ast.IncDecStmt:
		_, ok := b.(*ast.IncDecStmt)
		if ok {
			score += 1
		}
	case *ast.AssignStmt:
		_, ok := b.(*ast.AssignStmt)
		if ok {
			score += 1
		}
	case *ast.GoStmt:
		_, ok := b.(*ast.GoStmt)
		if ok {
			score += 1
		}
	case *ast.DeferStmt:
		_, ok := b.(*ast.DeferStmt)
		if ok {
			score += 1
		}
	case *ast.ReturnStmt:
		_, ok := b.(*ast.ReturnStmt)
		if ok {
			score += 1
		}
	case *ast.BranchStmt:
		_, ok := b.(*ast.BranchStmt)
		if ok {
			score += 1
		}
	case *ast.BlockStmt:
		_, ok := b.(*ast.BlockStmt)
		if ok {
			score += 1
		}
	case *ast.IfStmt:
		_, ok := b.(*ast.IfStmt)
		if ok {
			score += 1
		}
	case *ast.SwitchStmt:
		_, ok := b.(*ast.SwitchStmt)
		if ok {
			score += 1
		}
	case *ast.TypeSwitchStmt:
		_, ok := b.(*ast.TypeSwitchStmt)
		if ok {
			score += 1
		}
	case *ast.SelectStmt:
		_, ok := b.(*ast.SelectStmt)
		if ok {
			score += 1
		}
	case *ast.ForStmt:
		_, ok := b.(*ast.ForStmt)
		if ok {
			score += 1
		}
	case *ast.RangeStmt:
		_, ok := b.(*ast.RangeStmt)
		if ok {
			score += 1
		}
	}
	return
}
