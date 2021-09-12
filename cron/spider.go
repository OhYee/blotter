package cron

import (
	"sort"
	"sync"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/cron/spider"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/blotter/register"
)

func spiderSite(f friends.Friend, wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}

	friendName := f.Name
	friendURL := f.RSS

	var posts []friends.FriendPost

	retry := 0
	for retry = 0; retry < 5; retry++ {
		output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", friendName, friendURL, "retry", retry)
		posts = spider.Do(friendURL, retry)
		if len(posts) != 0 {
			if len(posts) > 5 {
				sort.SliceStable(posts, func(i, j int) bool { return posts[i].Time > posts[j].Time })
				posts = posts[:5]
			}
			break
		}
	}

	friends.SetFriendPosts(
		f.Link,
		posts,
	)
	output.DebugOutput.Println(posts)

	output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", friendName, friendURL, "Finished", retry)
}

func Spider() {
	output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider")
	defer output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", "Finished")

	spiderURLContext, ok := register.GetContext("spiderURL")
	spiderURL := ""
	if ok {
		switch spiderURLContext.(type) {
		case string:
			spiderURL = spiderURLContext.(string)
		}
	}

	if spiderURL == "" {
		wg := &sync.WaitGroup{}
		fs, _ := friends.GetFriends()
		for _, f := range fs {
			if f.RSS == "" {
				continue
			}
			wg.Add(1)
			go spiderSite(f, wg)
		}
		wg.Wait()
	} else {
		spiderSite(friends.Friend{
			Simple: friends.Simple{
				Name: "Test",
			},
			RSS: spiderURL,
		}, nil)
	}
}
