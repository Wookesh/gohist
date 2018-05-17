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
	"sync"
	"sync/atomic"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/gohist/objects"
	"github.com/wookesh/semaphore"
	"gitlab2.websensa.com/ma/websensa/logger"
	"gopkg.in/src-d/go-git.v4"
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
		return nil, fmt.Errorf("commitIterator is nil")
	}

	commitsData := make(map[string]*object.Commit)
	total := 0
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

	history := objects.NewHistory()

	last, first, graph := createGraph(commitsData, start, end)
	parentLocks := make(map[string]semaphore.Semaphore)
	for sha, node := range graph {
		total++
		if len(node.Parents) > 0 {
			s := semaphore.New(len(node.Parents))
			s.Empty()
			parentLocks[sha] = s
		} else {
			parentLocks[sha] = nil
		}
	}
	done := int32(0)
	queue := make(chan *Node, 1)
	queue <- first
	var wg sync.WaitGroup
	queued := make(map[string]bool)
	var m sync.Mutex
	for node := range queue {
		wg.Add(1)
		go func(node *Node) {
			atomic.AddInt32(&done, 1)
			logrus.Infoln("done:", atomic.LoadInt32(&done), "/", total)
			for _, child := range node.Children {
				defer func(child *Node) {
					parentLocks[child.SHA()].V()
				}(child)
			}
			defer wg.Done()
			for i := 0; i < len(node.Parents); i++ {
				parentLocks[node.SHA()].P()
			}

			files, err := node.Commit.Files()
			if err != nil {
				logrus.Fatalln(err)
			}
			var count int32
			var changed int32
			err = files.ForEach(func(f *object.File) error {
				if strings.Contains(f.Name, "vendor") || strings.Contains(f.Name, "Godeps") {
					// skip
					return nil
				}
				if !strings.HasSuffix(f.Name, ".go") || (strings.HasSuffix(f.Name, "_test.go") && !withTests) {
					return nil
				}
				logrus.Debugln("CreateHistory:", "\t", f.Name)
				rd, err := f.Blob.Reader()
				if err != nil {
					return err
				}
				body, err := ioutil.ReadAll(rd)
				if err != nil {
					logrus.Error("file.ForEach:", err)
					return err
				}
				functions, err := GetFunctions(string(body), f.Name, path.Dir(f.Name))
				if err != nil {
					logger.Warningln("CreateHistory:", "parse error:", err, f.Name)
					return nil
				}
				for funcID, funcDeclaration := range functions {
					added := history.Get(funcID).AddElement(funcDeclaration, node.Commit, body)
					if added {
						atomic.AddInt32(&changed, 1)
					}
					atomic.AddInt32(&count, 1)
				}
				return nil
			})
			if err != nil {
				logrus.Fatalln(err)
			}

			if changed > history.MaxChanged {
				history.MaxChanged = changed
			}

			atomic.AddInt32(&history.CommitsAnalyzed, 1)
			history.Mark(node.Commit.Author.When, int(count))
			history.CheckForDeleted(node.Commit)

			if node == last {
				close(queue)
			} else {
				for _, child := range node.Children {
					m.Lock()
					_, ok := queued[child.SHA()]
					if !ok {
						queued[child.SHA()] = true
						queue <- child
					}
					m.Unlock()
				}
			}
		}(node)
	}

	wg.Wait()

	for _, f := range history.Data {
		f.PostProcess()
	}

	return history, nil
}

type Node struct {
	Commit   *object.Commit
	Children []*Node
	Parents  []*Node
}

func (n *Node) SHA() string {
	return n.Commit.Hash.String()
}

func createGraph(commits map[string]*object.Commit, start, end string) (first, last *Node, graph map[string]*Node) {
	graph = make(map[string]*Node)
	for k, elem := range commits {
		graph[k] = &Node{Commit: elem}
	}
	for _, elem := range graph {
		for _, hash := range elem.Commit.ParentHashes {
			parent, ok := graph[hash.String()]
			if ok {
				parent.Children = append(parent.Children, elem)
				elem.Parents = append(elem.Parents, parent)
			}
		}
	}
	first = graph[start]
	if elem, ok := graph[end]; end != "" && ok {
		last = elem
	} else {
		i := first
		for len(i.Parents) > 0 {
			i = i.Parents[0]
		}
		last = i
	}

	counts := make(map[string]int)
	var queue []*Node
	visited := make(map[string]bool)
	queue = append(queue, first)
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		_, ok := visited[elem.Commit.Hash.String()]
		if ok {
			continue
		}
		visited[elem.Commit.Hash.String()] = true
		counts[elem.Commit.Hash.String()] += 1
		for _, parent := range elem.Parents {
			queue = append(queue, parent)
		}
	}

	visited = make(map[string]bool)
	queue = append(queue, last)
	for len(queue) > 0 {
		elem := queue[0]
		queue = queue[1:]
		_, ok := visited[elem.Commit.Hash.String()]
		if ok {
			continue
		}
		visited[elem.Commit.Hash.String()] = true
		counts[elem.Commit.Hash.String()] += 1
		for _, child := range elem.Children {
			queue = append(queue, child)
		}
	}

	for hash, node := range graph {
		count := counts[hash]
		if count < 2 {
			delete(graph, hash)
		}
		var clearParents []*Node
		for _, parent := range node.Parents {
			if count, ok := counts[parent.SHA()]; ok && count == 2 {
				clearParents = append(clearParents, parent)
			}
		}
		node.Parents = clearParents
		var clearChildren []*Node
		for _, parent := range node.Children {
			if count, ok := counts[parent.SHA()]; ok && count == 2 {
				clearChildren = append(clearChildren, parent)
			}
		}
		node.Children = clearChildren
	}

	return first, last, graph
}

func GetFunctions(src, fileName, pack string) (map[string]*ast.FuncDecl, error) {
	fileSet := token.NewFileSet()
	f, err := parser.ParseFile(fileSet, "", src, parser.AllErrors)
	if err != nil {
		return nil, err
	}
	functions := make(map[string]*ast.FuncDecl)
	//variables := make(map[string]*objects.Variable)
	for _, decl := range f.Decls {
		if function, ok := decl.(*ast.FuncDecl); ok {
			functions[pack+"."+createSignature(function, fileName)] = function
		}
		//if v, ok := decl.(*ast.GenDecl); ok {
		//	switch v.Tok {
		//	case token.VAR:
		//		gatherVariables(v, variables)
		//	case token.TYPE:
		//	case token.IMPORT:
		//	case token.CONST:
		//		gatherVariables(v, variables)
		//	}
		//}
	}
	//_ = variables
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

func createSignature(f *ast.FuncDecl, fileName string) (signature string) {
	if f == nil {
		return
	}
	name := f.Name.Name
	if name == "init" {
		name = fmt.Sprintf("%s[%s]", name, fileName)
	}
	if f.Recv != nil {
		var recv []string
		for _, param := range f.Recv.List {
			if len(param.Names) == 0 {
				recv = append(recv, getType(param.Type))
			}
			for range param.Names {
				recv = append(recv, getType(param.Type))
			}
		}

		return strings.Join(recv, ",") + "." + name
	}
	return name
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
