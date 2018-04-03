package main

import (
	"flag"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/gohist/collector"
	"github.com/wookesh/gohist/ui"
)

var (
	projectPath = flag.String("path", "", "path to repo")
	port        = flag.String("port", "8000", "port for web server")
	start       = flag.String("start", "master", "newest commit to parse")
	end         = flag.String("end", "", "latest commit to parse")
)

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	flag.Parse()

	if *projectPath == "" {
		flag.PrintDefaults()
		return
	}

	absProjectPath, err := filepath.Abs(*projectPath)
	if err != nil {
		logrus.Fatalln(err)
	} else {
		*projectPath = absProjectPath
	}

	history, err := collector.CreateHistory(*projectPath, *start, *end, false)
	if err != nil {
		panic(err)
	}
	_ = history
	split := strings.Split(*projectPath, "/src/")
	var repoName string
	if len(split) >= 2 {
		repoName = split[1]
	} else {
		repoName = *projectPath
	}
	ui.Run(history, repoName, *port)
}
