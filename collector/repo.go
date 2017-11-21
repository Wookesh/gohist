package collector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/wookesh/mgr/objects"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func CreateHistory(path string) *objects.History {
	repo, _ := git.PlainOpen(path)
	iter, _ := repo.CommitObjects()
	iter.ForEach(func(commit *object.Commit) error {
		fmt.Println("##################################")
		fmt.Println(commit.Hash)
		fmt.Println("##################################")
		files, _ := commit.Files()
		files.ForEach(func(f *object.File) error {
			if strings.HasSuffix(f.Name, ".go") {
				rd, _ := f.Blob.Reader()
				body, err := ioutil.ReadAll(rd)
				if err != nil {
					fmt.Println(err)
				}
				GetFunctions(string(body))
			}
			return nil
		})
		return nil
	})
	return nil
}

func GetFunctions(src string) (map[string]*ast.FuncDecl, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	var functions map[string]*ast.FuncDecl
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			functions[function.Name.Name] = function
		}
	}
	return functions, nil
}
