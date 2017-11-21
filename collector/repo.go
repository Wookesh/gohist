package collector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/wookesh/gohist/objects"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func CreateHistory(path string) (*objects.History, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}
	iter, _ := repo.CommitObjects()
	data := make(map[*object.Commit]map[string]*ast.FuncDecl)
	iter.ForEach(func(commit *object.Commit) error {
		files, _ := commit.Files()
		data[commit] = make(map[string]*ast.FuncDecl)
		files.ForEach(func(f *object.File) error {
			if strings.HasSuffix(f.Name, ".go") {
				rd, err := f.Blob.Reader()
				if err != nil {
					return err
				}
				body, err := ioutil.ReadAll(rd)
				if err != nil {
					fmt.Println(err)
				}
				functions, err := GetFunctions(string(body))
				if err != nil {
					return err
				}
				for k, v := range functions {
					data[commit][k] = v
				}
			}
			return nil
		})
		return nil
	})
	return nil, nil
}

func GetFunctions(src string) (map[string]*ast.FuncDecl, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	functions := make(map[string]*ast.FuncDecl)
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			functions[function.Name.Name] = function
		}
	}
	return functions, nil
}
