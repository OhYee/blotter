package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/cron"
	"github.com/OhYee/blotter/http"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/blotter/utils/initial"
	"github.com/OhYee/blotter/utils/reaper"
)

var (
	addr   = "127.0.0.1:50000"
	prefix = "/api/"
	tool   = ""
	url    = ""
)

var (
	_version string
)

//go:generate /bin/bash ./generate.bash

func parseFlags() {
	flag.StringVar(&addr, "address", "127.0.0.1:50000", "listen address")
	flag.StringVar(&prefix, "prefix", "/api/", "api url prefix")
	flag.StringVar(&url, "url", "", "spider url")

	keys := make([]string, len(cron.CronMap))
	pos := 0
	for k, _ := range cron.CronMap {
		keys[pos] = k
		pos++
	}
	flag.StringVar(&tool, "tool", "", fmt.Sprintf("call tools(%s)", strings.Join(keys, ",")))

	flag.Parse()
}

func main() {
	go reaper.Reap()

	parseFlags()
	register.SetContext("version", _version)
	register.SetContext("spiderURL", url)

	if tool != "" {
		f, e := cron.CronMap[tool]
		if !e {
			output.ErrOutput.Printf("No tool named %s\n", tool)
			os.Exit(1)
		}
		f()
		os.Exit(0)
		return
	}

	cron.Start()
	defer cron.Stop()

	api.Register()
	initial.Run()
	// queue.Register().Register("extensions/queue")
	// register.DebugApiMap()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
