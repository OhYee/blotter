package cron

import (
	"sort"
	"sync"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/api/pkg/variable"
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

func sortFriends() {
	fs, err := friends.GetFriends()
	if err != nil {
		output.Err(err)
		return
	}

	sort.SliceStable(fs, func(i, j int) bool {
		if fs[i].Ex != fs[j].Ex {
			return fs[i].Ex
		}
		var iTime, jTime int64
		if len(fs[i].Posts) > 0 {
			iTime = fs[i].Posts[0].Time
		}
		if len(fs[j].Posts) > 0 {
			jTime = fs[j].Posts[0].Time
		}
		return !(iTime < jTime)
	})

	// move the link is root to the first
	root := ""
	vMap, err := variable.Get("root")
	if err == nil {
		root, _ = vMap.GetString("root")
	}
	if root != "" {
		for i, f := range fs {
			if f.Link == root {
				temp := []friends.Friend{f}
				if i > 0 {
					temp = append(temp, fs[:i]...)
				}
				if i+1 < len(fs) {
					temp = append(temp, fs[i+1:]...)
				}
				fs = temp
				break
			}
		}
	}

	err = friends.SetFriends(fs)
	if err != nil {
		output.Err(err)
	}
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
		// fs, _ := friends.GetFriends()
		// for _, f := range fs {
		// 	if f.RSS == "" {
		// 		continue
		// 	}
		// 	wg.Add(1)
		// 	func(f friends.Friend, wg *sync.WaitGroup) {
		// 		pool.Do(func() {
		// 			spiderSite(f, wg)
		// 		})
		// 	}(f, wg)
		// }
		wg.Wait()
		sortFriends()
	} else {
		spiderSite(friends.Friend{
			Simple: friends.Simple{
				Name: "Test",
			},
			RSS: spiderURL,
		}, nil)
	}
}
