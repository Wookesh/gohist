package collector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"strings"

	"github.com/wookesh/gohist/objects"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type node struct {
	parents  []*node
	commit   *object.Commit
	children []*node
}

func CreateHistory(repoPath string, start, end string, withTests bool) (*objects.History, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	commitIterator, _ := repo.CommitObjects()
	_, err = repo.CommitObject(plumbing.NewHash(start))
	if err != nil {
		ref, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+start), false)
		if err != nil {
			return nil, err
		}
		start = ref.Hash().String()
	}

	commitsData := make(map[string]*object.Commit)
	commitIterator.ForEach(func(commit *object.Commit) error {
		commitsData[commit.Hash.String()] = commit
		return nil
	})
	historyLine := createHistoryLine(commitsData, start, end)
	history := objects.NewHistory()

	for _, commit := range historyLine {
		fmt.Println(commit.Hash)
		files, err := commit.Files()
		if err != nil {
			return nil, err
		}
		files.ForEach(func(f *object.File) error {
			if strings.HasSuffix(f.Name, ".go") && (!strings.HasSuffix(f.Name, "_test.go") || withTests) {
				fmt.Println("\t", f.Name)
				rd, err := f.Blob.Reader()
				if err != nil {
					return err
				}
				body, err := ioutil.ReadAll(rd)
				if err != nil {
					fmt.Println(err)
				}
				functions, err := GetFunctions(string(body), path.Dir(f.Name))
				if err != nil {
					return err
				}
				for funcID, funcDecl := range functions {
					funcHistory, ok := history.Data[funcID]
					if !ok {
						funcHistory = &objects.FunctionHistory{}
						history.Data[funcID] = funcHistory
					}
					text := string(body[funcDecl.Pos()-1 : funcDecl.End()-1])
					funcHistory.History = append(funcHistory.History,
						&objects.HistoryElement{
							Func:   funcDecl,
							Commit: commit,
							Text:   text,
							Offset: int(funcDecl.Pos()),
						})
				}
			}
			return nil
		})
	}
	return history, nil
}

func createHistoryLine(commits map[string]*object.Commit, start, end string) (history []*object.Commit) {
	rootNode, ok := commits[start]
	if !ok {
		return
	}
	var queue [][]*object.Commit
	queue = append(queue, []*object.Commit{rootNode})
	visited := make(map[string]bool)
	for len(queue) > 0 {
		filePath := queue[0]
		queue = queue[1:]
		history = filePath
		elem := filePath[len(filePath)-1]
		if _, ok := visited[elem.Hash.String()]; ok {
			continue
		}
		visited[elem.Hash.String()] = true

		for _, hash := range elem.ParentHashes {
			commit := commits[hash.String()]
			if commit != nil {
				if hash.String() == end {
					return append(filePath, commit)
				}
				newPath := make([]*object.Commit, len(filePath))
				copy(newPath, filePath)
				newPath = append(newPath, commit)
				queue = append(queue, newPath)
			}
		}
	}
	for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
		history[i], history[j] = history[j], history[i]
	}
	return
}

func GetFunctions(src, pack string) (map[string]*ast.FuncDecl, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	functions := make(map[string]*ast.FuncDecl)
	variables := make(map[string]*objects.Variable)
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			functions[pack+"."+function.Name.Name] = function
		}
		if v, ok := decl.(*ast.GenDecl); ok {
			switch v.Tok {
			case token.VAR:
				for _, spec := range v.Specs {
					value, _ := spec.(*ast.ValueSpec)
					for i := 0; i < len(value.Names); i++ {
						v := &objects.Variable{Name: value.Names[i], Type: value.Type}
						if len(value.Values) > 0 {
							v.Expr = value.Values[i]
						}
						variables[v.Name.Name] = v
					}
				}
			case token.TYPE:
			case token.IMPORT:
			case token.CONST:
			}
		}
	}
	_ = variables
	return functions, nil
}
