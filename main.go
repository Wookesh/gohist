package main

import (
	"flag"
	"net/http"
	_ "net/http/pprof"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/wookesh/gohist/collector"
	"github.com/wookesh/gohist/ui"
)

var (
	projectPath = flag.String("path", ".", "path to repo")
	port        = flag.String("port", "8000", "port for web server")
	start       = flag.String("start", "master", "newest commit to parse")
	end         = flag.String("end", "", "latest commit to parse")
	debug       = flag.Bool("debug", false, "Run debug mode")
	simple      = flag.Bool("simple_diff", false, "Create graph using standard diff")
)

func main() {
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

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

	history, err := collector.CreateHistory(*projectPath, *start, *end, false, *simple)
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
