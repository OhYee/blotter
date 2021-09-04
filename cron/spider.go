package cron

import (
	"sync"
	"time"

	"strings"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/cron/spider"
	"github.com/OhYee/blotter/output"
)

func Spider() {
	output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider")
	defer output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", "Finished")

	wg := &sync.WaitGroup{}

	fs, _ := friends.GetFriends()
	for _, f := range fs {
		wg.Add(1)
		go func(f friends.Friend) {
			defer wg.Done()

			friendName := f.Name
			friendURL := f.RSS

			var posts []friends.FriendPost

			retry := 0
			for retry = 0; retry < 5; retry++ {
				output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", friendName, friendURL, "retry", retry)
				if friendURL == "" ||
					strings.Index(friendURL, "rss") != -1 ||
					strings.Index(friendURL, "atom") != -1 ||
					strings.Index(friendURL, "feed") != -1 ||
					strings.Index(friendURL, "xml") != -1 {
					posts = spider.ReadRSS(friendURL)
				} else {
					posts = spider.ReadHTML(friendURL)
				}
				if len(posts) != 0 {
					break
				}
			}

			friends.SetFriendPosts(
				f.Link,
				posts,
			)

			output.LogOutput.Println(time.Now().Format("2006-01-02 15:04:05"), "Spider", friendName, friendURL, "Finished", retry)
		}(f)

	}
	wg.Wait()
}
