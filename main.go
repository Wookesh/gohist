package main

import (
	"flag"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/wookesh/gohist/collector"
	"github.com/wookesh/gohist/ui"
)

var (
	projectPath = flag.String("path", "/home/wookesh/GoProjects/src/github.com/wookesh/gohist", "")
	start       = flag.String("start", "master", "")
	end         = flag.String("end", "", "")
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	flag.Parse()

	history, err := collector.CreateHistory(*projectPath, *start, *end, false)
	if err != nil {
		panic(err)
	}
	split := strings.Split(*projectPath, "/src/")
	var repoName string
	if len(split) >= 2 {
		repoName = split[1]
	} else {
		repoName = *projectPath
	}

	ui.Run(history, repoName)
}
