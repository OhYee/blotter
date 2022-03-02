package spider

import (
	"regexp"
	"sort"
	"strings"

	"github.com/OhYee/blotter/api/pkg/friends"
)

var ignoreIllegalCharacter = regexp.MustCompile(`[\x00-\x09\x0b\x0c\x0e-\x1f]`)

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

	// ignore illegal character suchj as BS
	content = ignoreIllegalCharacter.ReplaceAllString(content, "")

	if strings.Contains(link, "rss") ||
		strings.Contains(link, "atom") ||
		strings.Contains(link, "feed") ||
		strings.Contains(link, "xml") {
		posts = readRSS(link, content)
	} else if isJSON || strings.Contains(link, "json") {
		posts = readJSON(link, content)
	} else {
		posts = readHTML(link, content)
	}

	sort.Slice(posts, func(i, j int) bool { return posts[i].Time > posts[j].Time })
	return posts
}
