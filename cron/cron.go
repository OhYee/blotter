package cron

import (
	"github.com/OhYee/blotter/output"
	"github.com/robfig/cron"
)

var c *cron.Cron = nil

var CronMap = map[string]func(){
	"spider":    Spider,
	"backup":    Backup,
	"baidupush": BaiduPush,
}

func checkError(err error) error {
	if err != nil {
		output.ErrOutput.Println(err)
	}
	return err
}

func Start() {
	if c != nil {
		c.Stop()
	}

	c = cron.New() //精确到秒

	// 三点启动爬虫任务
	checkError(c.AddFunc("0 0 3 * * ?", Spider))

	// 三点半启动定时备份任务
	checkError(c.AddFunc("0 30 3 * * ?", Backup))

	// 每小时启动定时百度推送任务
	checkError(c.AddFunc("0 0 * * * ?", BaiduPush))

	c.Start()
}

func Stop() {
	if c != nil {
		c.Stop()
	}
}
