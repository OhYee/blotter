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
)

var (
	addr   = "127.0.0.1:50000"
	prefix = "/api/"
	tool   = ""
)

var (
	_version string
)

//go:generate /bin/bash ./generate.bash

func parseFlags() {
	flag.StringVar(&addr, "address", "127.0.0.1:50000", "listen address")
	flag.StringVar(&prefix, "prefix", "/api/", "api url prefix")

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
	parseFlags()
	if tool != "" {
		f, e := cron.CronMap[tool]
		if !e {
			output.ErrOutput.Printf("No tool named %s\n", tool)
			os.Exit(1)
		}
		f()
		os.Exit(0)
	}

	register.SetContext("version", _version)

	cron.Start()
	defer cron.Stop()

	api.Register()
	// queue.Register().Register("extensions/queue")
	// register.DebugApiMap()
	output.Log("Server will start at http://%s", addr)
	if err := http.Server(addr, prefix); err != nil {
		output.Err(err)
	}
}
