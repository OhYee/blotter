package spider

import (
	"bytes"
	"strings"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/output"
	"github.com/mmcdole/gofeed"
	"golang.org/x/net/html"
)

func findTitle(n *html.Node, arr []*html.Node) {
	t := n
	for t != nil {
		if t.Type == html.ElementNode && t.Data == "title" {
			arr = append(arr, t)
		}
		findTitle(t.FirstChild, arr)
		t = t.NextSibling
	}
}

// htmlParser 以 html 格式解析 RSS
func htmlParser(host, content string) (feed *gofeed.Feed, err error) {
	root, err := html.Parse(bytes.NewReader([]byte(content)))
	if err != nil {
		return
	}
	titles := make([]*html.Node, 0)
	findTitle(root, titles)

	feed = &gofeed.Feed{}
	feed.Items = make([]*gofeed.Item, len(titles))
	for i, title := range titles {
		pointer := title
		for pointer.PrevSibling != nil {
			pointer = pointer.PrevSibling
		}

		var ts *time.Time = nil
		link := ""

		for pointer != nil {
			// 判断是否是日期
			value := pointer.Data
			if ts == nil {
				if tmp := parseTime(value); tmp != nil {
					ts = tmp
				}
			}

			// 判断是否是 link
			if link == "" {
				if strings.HasPrefix(value, "https://") ||
					strings.HasPrefix(value, "http://") ||
					strings.HasPrefix(value, "//") {
					link = value
				}
			}
		}
		feed.Items[i] = &gofeed.Item{
			Title:         title.Data,
			Link:          link,
			UpdatedParsed: ts,
		}
	}
	return
}

func readRSS(link, content string) (posts []friends.FriendPost) {
	output.DebugOutput.Println(link, "readRSS")
	posts = make([]friends.FriendPost, 0, 5)

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(content)
	if err != nil {
		output.ErrOutput.Println(link, err)
		// gofeed 解析失败，换用 html 形式解析
		if feed, err = htmlParser(link, content); err != nil {
			output.ErrOutput.Println(link, err)
			return
		}
	}

	for _, p := range feed.Items {
		updateTime := p.UpdatedParsed
		if updateTime == nil {
			updateTime = p.PublishedParsed
		}
		if updateTime == nil {
			updateTime = parseTime(p.Updated)
		}
		if updateTime == nil {
			updateTime = parseTime(p.Published)
		}
		// output.DebugOutput.Println(p)
		posts = append(posts, friends.FriendPost{
			Title: p.Title,
			Link:  p.Link,
			Time:  toUnix(updateTime),
		})
	}

	return posts
}
