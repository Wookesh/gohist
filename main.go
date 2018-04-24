package main

import (
	"flag"
	_ "net/http/pprof"
	"path/filepath"
	"strings"

	"net/http"

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

	go func() { http.ListenAndServe(":6060", nil) }()

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

	split := strings.Split(*projectPath, "/src/")
	var repoName string
	if len(split) >= 2 {
		repoName = split[1]
	} else {
		repoName = *projectPath
	}
	ui.Run(history, repoName, *port)
}
