package main

import (
	"flag"
	"fmt"

	"github.com/fatih/color"
	"github.com/wookesh/gohist/collector"
	"github.com/wookesh/gohist/diff"
)

var (
	projectPath = flag.String("path", "C:/Go/projects/src/github.com/wookesh/distributed", "")
)

func main() {
	flag.Parse()

	history, err := collector.CreateHistory(*projectPath)
	if err != nil {
		panic(err)
	}
	fmt.Print(history)

	//fset := token.NewFileSet()
	//
	//pkgs, err := parser.ParseDir(fset, *projectPath, nil, parser.AllErrors)
	//if err != nil {
	//	panic(err)
	//}
	//ast.Print(fset, pkgs)
	//for name, pkg := range pkgs {
	//	log.Println("Package:", name)
	//	for fileName, file := range pkg.Files {
	//		log.Println("File:", fileName)
	//		text, err := ioutil.ReadFile(path.Join(fileName))
	//		if err != nil {
	//			panic(err)
	//		}
	//		for _, decl := range file.Decls {
	//			if f, ok := decl.(*ast.FuncDecl); ok {
	//				log.Println("Func:", f.Name)
	//				//log.Println(string(text[f.Pos()-1:f.End()]))
	//				coloring := diff.Diff(f, f, diff.ModeBoth)
	//				log.Println(coloring)
	//				//log.Print(len(text), f.End()-1)
	//				//if len(text) == int(f.End()-1) {
	//				//	 last function, no ENDLINE
	//				//PrintWithColor(coloring, string(text),int(f.Pos()-1), int(f.End()-2))
	//				//} else {
	//				PrintWithColor(coloring, string(text), int(f.Pos()-1), int(f.End()-1))
	//				//}
	//			}
	//		}
	//	}
	//}
}

func PrintWithColor(coloring diff.Coloring, text string, start, end int) {
	defer color.Unset()
	if len(coloring) == 0 {
		fmt.Print(text[start:end])
		return
	}
	current := 0
	var hasColoring bool
	for i := start; int(i) <= end; i++ {
		if !hasColoring && int(coloring[current].Pos)-1 <= i {
			hasColoring = true
			setColor(coloring[current].Color)
		}

		if hasColoring && int(coloring[current].End)-1 < i {
			color.Unset()
			if current <= len(coloring)-2 {
				current++
			}
			if int(coloring[current].Pos)-1 >= i {
				setColor(coloring[current].Color)
			} else {
				hasColoring = false
			}
		}

		fmt.Print(string(text[i]))
	}
}

func setColor(c diff.Color) {
	switch c {
	case diff.ColorSame:
		color.Set(color.FgWhite)
	case diff.ColorNew:
		color.Set(color.FgGreen)
	case diff.ColorRemoved:
		color.Set(color.FgRed)
	}
}
