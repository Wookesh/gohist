package collector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"path"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wookesh/gohist/diff"
	"github.com/wookesh/gohist/objects"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func CreateHistory(repoPath string, start, end string, withTests bool) (*objects.History, error) {
	logrus.Debugln("CreateHistory:", repoPath)
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}
	commitIterator, err := repo.CommitObjects()
	if err != nil {
		return nil, err
	}

	if commitIterator == nil {
		return nil, fmt.Errorf("commitIterator is nil!")
	}

	commitsData := make(map[string]*object.Commit)
	commitIterator.ForEach(func(commit *object.Commit) error {
		if commit == nil {
			panic("commit is nil")
			return fmt.Errorf("commit is nil")
		}
		commitsData[commit.Hash.String()] = commit
		return nil
	})

	_, err = repo.CommitObject(plumbing.NewHash(start))
	if err != nil {
		ref, err := repo.Reference(plumbing.ReferenceName("refs/heads/"+start), false)
		if err != nil {
			return nil, err
		}
		start = ref.Hash().String()
	}

	historyLine := createHistoryLine(commitsData, start, end)
	history := objects.NewHistory()

	for _, commit := range historyLine {
		logrus.Debugln("CreateHistory:", commit.Hash)
		files, err := commit.Files()
		if err != nil {
			return nil, err
		}
		added := make(map[string]bool)
		files.ForEach(func(f *object.File) error {
			if strings.HasSuffix(f.Name, ".go") && (!strings.HasSuffix(f.Name, "_test.go") || withTests) {
				logrus.Debugln("CreateHistory:", "\t", f.Name)
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
					added[funcID] = true
					if len(funcHistory.History) > 0 {
						last := len(funcHistory.History) - 1
						if diff.IsSame(funcHistory.History[last].Func, funcDecl) {
							continue
						}
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
		for funcID, funcHistory := range history.Data {
			if _, ok := added[funcID]; !ok {
				if funcHistory.History[len(funcHistory.History)-1].Func != nil {
					funcHistory.History = append(funcHistory.History, &objects.HistoryElement{
						Func:   nil,
						Commit: commit,
						Text:   "",
						Offset: 0,
					})
				}
			}
		}
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
					history = append(filePath, commit)
					for i, j := 0, len(history)-1; i < j; i, j = i+1, j-1 {
						history[i], history[j] = history[j], history[i]
					}
					return history
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
			functions[pack+"."+createSignature(function)] = function
		}
		if v, ok := decl.(*ast.GenDecl); ok {
			switch v.Tok {
			case token.VAR:
				gatherVariables(v, variables)
			case token.TYPE:
			case token.IMPORT:
			case token.CONST:
				gatherVariables(v, variables)
			}
		}
	}
	_ = variables
	return functions, nil
}

func gatherVariables(v *ast.GenDecl, variables map[string]*objects.Variable) {
	for _, spec := range v.Specs {
		value, _ := spec.(*ast.ValueSpec)
		for i := 0; i < len(value.Names); i++ {
			v := &objects.Variable{Name: value.Names[i], Type: value.Type}
			if len(value.Values) > i {
				v.Expr = value.Values[i]
			}
			variables[v.Name.Name] = v
		}
	}
}

func createSignature(f *ast.FuncDecl) (signature string) {
	if f == nil {
		return
	}
	//if f.Type.Params != nil {
	//	for _, param := range f.Type.Params.List {
	//		for _, name := range param.Names {
	//			logrus.Infoln("param:", name.Name, getType(param.Type))
	//		}
	//	}
	//}
	//
	//if f.Type.Results != nil {
	//	for _, param := range f.Type.Results.List {
	//		for _, name := range param.Names {
	//			logrus.Infoln("result:", name.Name, getType(param.Type))
	//		}
	//	}
	//}
	if f.Recv != nil {
		var recv []string
		for _, param := range f.Recv.List {
			for range param.Names {
				recv = append(recv, getType(param.Type))
			}
		}
		return strings.Join(recv, ",") + "." + f.Name.Name
	}
	return f.Name.Name
}

func getType(x ast.Node) string {
	if x == nil {
		return ""
	}
	switch t := x.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.SelectorExpr:
		return getType(t.X) + "." + t.Sel.Name
	case *ast.StarExpr:
		return getType(t.X)
	case *ast.ArrayType:
		return "[" + getType(t.Len) + "]" + getType(t.Elt)
	case *ast.MapType:
		return "map[" + getType(t.Key) + "]" + getType(t.Value)
	case *ast.InterfaceType:
		if len(t.Methods.List) == 0 {
			return "interface{}"
		} else {
			panic(reflect.TypeOf(t))
			return ""
		}
	default:
		panic(reflect.TypeOf(t))
	}
}
