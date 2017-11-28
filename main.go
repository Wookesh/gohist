package main

import (
	"flag"

	"github.com/wookesh/gohist/collector"
	"github.com/wookesh/gohist/ui"
)

var (
	projectPath = flag.String("path", "/home/wookesh/GoProjects/src/github.com/wookesh/gohist", "")
	start       = flag.String("start", "master", "")
	end         = flag.String("end", "", "")
)

func main() {
	flag.Parse()

	history, err := collector.CreateHistory(*projectPath, *start, *end, false)
	if err != nil {
		panic(err)
	}
	ui.Run(history)
}
