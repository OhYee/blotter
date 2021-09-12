package spider

import (
	"sort"
	"strings"

	"github.com/OhYee/blotter/api/pkg/friends"
)

func Do(link string, retry int) (posts []friends.FriendPost) {
	posts = make([]friends.FriendPost, 0)
	if link == "" {
		return
	}
	content := ""
	isJSON := false
	if retry%2 == 0 {
		content, isJSON = getHTML(link)
	} else {
		content = getHTMLWithJS(link)
	}
	if len(content) == 0 {
		return
	}

	if strings.Index(link, "rss") != -1 ||
		strings.Index(link, "atom") != -1 ||
		strings.Index(link, "feed") != -1 ||
		strings.Index(link, "xml") != -1 {
		posts = readRSS(link, content)
	} else if isJSON || strings.Index(link, "json") != -1 {
		posts = readJSON(link, content)
	} else {
		posts = readHTML(link, content)
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].Time > posts[j].Time })
	return posts
}
