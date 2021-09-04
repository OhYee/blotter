package spider

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/output"
	"github.com/mmcdole/gofeed"
)

func ReadRSS(u string) (posts []friends.FriendPost) {
	// output.DebugOutput.Println(u)
	posts = make([]friends.FriendPost, 0, 5)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: Timeout,
	}
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}

	req.Header.Set("User-Agent", UserAgent)

	resp, err := client.Do(req)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		output.ErrOutput.Println(u, err)
		return
	}

	// output.DebugOutput.Println(string(b))

	fp := gofeed.NewParser()
	feed, err := fp.ParseString(string(b))
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
