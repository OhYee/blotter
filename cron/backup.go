package cron

import (
	"time"

	"github.com/OhYee/blotter/output"
)

func Backup() {
	output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Backup")
}
