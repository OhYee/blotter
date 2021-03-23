package post

import (
	"sync"

	"github.com/OhYee/blotter/output"
	"github.com/yanyiwu/gojieba"
)

// 使用 sync.Once 延迟初始化 jieba
// 确保只会初始化一次

var jieba *gojieba.Jieba
var jiebaOnce = sync.Once{}

func initJieba() {
	jiebaOnce.Do(
		func() {
			output.Log("Initial Jieba")
			jieba = gojieba.NewJieba()
			output.Log("Initial Jieba finished")
		},
	)
}

func getJieba() *gojieba.Jieba {
	if jieba == nil {
		initJieba()
	}
	return jieba
}

func init() {
	go initJieba()
}
