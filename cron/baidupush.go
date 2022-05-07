package cron

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/output"
)

const url = "http://data.zz.baidu.com/urls?site=%s&token=%s"

func BaiduPush() {
	output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Baidu push")

	// 读入配置项
	variables, err := variable.Get("root", "baidupush")
	if err != nil {
		output.Err(err)
		return
	}
	root, _ := variables.GetString("root")
	token, _ := variables.GetString("baidupush")
	if root == "" || token == "" {
		output.LogOutput.Println("No need for ")
		return
	}

	// 生成链接文件
	buf := bytes.NewBufferString("")
	_, posts, err := post.GetCardPosts(0, 0, []string{}, []string{}, "", 0, "", []string{}, true)
	if err != nil {
		return
	}
	for _, post := range posts {
		buf.WriteString(fmt.Sprintf("%s/post/%s\n", root, post.URL))
	}

	// 发送请求
	c := http.Client{}
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(url, root, token),
		buf,
	)
	if err != nil {
		output.Err(err)
		return
	}
	req.Header.Set("Content-Type", "text/plain")
	resp, err := c.Do(req)
	output.Log(resp.Status, err)
	return
}
