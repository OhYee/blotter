package spider

import (
	"bytes"

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

func htmlParser(content string) (feed *gofeed.Feed, err error) {
	root, err := html.Parse(bytes.NewReader([]byte(content)))
	if err != nil {
		return
	}
	titles := make([]*html.Node, 0)
	findTitle(root, titles)

	feed = &gofeed.Feed{}
	feed.Items = make([]*gofeed.Item, len(titles))
	for i, title := range titles {
		feed.Items[i] = &gofeed.Item{
			Title:     title.Data,
			Link:      "",
			Updated:   "",
			Published: "",
		}
	}
	return
}

func ReadRSS(u string, retry int) (posts []friends.FriendPost) {
	// output.DebugOutput.Println(u)
	posts = make([]friends.FriendPost, 0, 5)

	content := ""
	if retry%2 == 0 {
		content = getHTML(u)
	} else {
		content = getHTMLWithJS(u)
	}
	if len(content) == 0 {
		return
	}

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(content)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
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

	if len(posts) > 5 {
		return posts[:5]
	}
	return posts
}
