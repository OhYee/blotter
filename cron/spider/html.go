package spider

import (
	"bytes"
	"strings"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"golang.org/x/net/html"

	"net/url"
)

func truncate(s string) string {
	if len(s) > 20 {
		return s[:20] + "..."
	}
	return s
}

func dfs(curr *html.Node, xpath []string, m map[string][]*html.Node) {
	if curr == nil {
		return
	}

	if curr.DataAtom.String() == "a" {
		x := strings.Join(xpath, "|")
		if _, ok := m[x]; !ok {
			m[x] = make([]*html.Node, 0)
		}
		m[x] = append(m[x], curr)
		// fmt.Printf("%s %s %s %s\n", strings.Join(xpath, "|"), curr.Attr, curr.Namespace, curr.Data)
	}

	child := curr.FirstChild
	for child != nil {
		dfs(child, append(xpath, curr.DataAtom.String()), m)
		child = child.NextSibling
	}
}

func elementInnterText(curr *html.Node) []string {
	if curr == nil {
		return []string{}
	}
	if curr.Type == 1 {
		data := strings.TrimSpace(curr.Data)
		if data == "" {
			return []string{}
		}
		return []string{data}
	}
	s := []string{}
	child := curr.FirstChild

	for child != nil {
		s = append(s, elementInnterText(child)...)
		child = child.NextSibling
	}
	return s
}

func elementHref(curr *html.Node) string {
	if curr == nil {
		return ""
	}
	for _, attr := range curr.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}

func CheckPost(v []*html.Node) bool {
	if len(v) < 2 {
		return false
	}
	sumTitle := 0.0
	hosts := make(map[string]struct{})
	for _, item := range v {
		for _, s := range elementInnterText(item) {
			sumTitle += float64(len(s))
		}

		href := elementHref(item)

		url, err := url.Parse(href)
		if err != nil {
			return false
		}
		hosts[url.Host] = struct{}{}
		if url.Path == "/" {
			return false
		}

	}
	if sumTitle/float64(len(v)) < 5*2 {
		return false
	}
	if len(hosts) > 1 {
		return false
	}
	return true
}

func elementFindTime(node *html.Node) *time.Time {
	t := node
	parentCount := 0
	for parentCount < 5 && t != nil {
		ss := elementInnterText(t)
		for _, s := range ss {
			ts := parseTime(s)
			if ts != nil {
				return ts
			}
		}

		parentCount++
		t = t.Parent
	}
	return nil
}

func ReadHTML(u string, retry int) []friends.FriendPost {
	hostURL, _ := url.Parse(u)

	content := ""
	if retry%2 == 0 {
		content = getHTML(u)
	} else {
		content = getHTMLWithJS(u)
	}
	// output.DebugOutput.Println(c)

	doc, _ := html.Parse(bytes.NewBufferString(content))
	m := make(map[string][]*html.Node)
	dfs(doc, []string{}, m)

	posts := make([]friends.FriendPost, 0, 5)

	var postList []*html.Node
	var maxPostListLength = 0
	for _, v := range m {
		if CheckPost(v) && len(v) > maxPostListLength {
			postList = v
			maxPostListLength = len(v)
		}
	}

	for _, item := range postList {
		u, _ := url.Parse(elementHref(item))
		if u.Host == "" {
			u.Scheme = hostURL.Scheme
			u.Host = hostURL.Host
		}

		titles := elementInnterText(item)

		// output.DebugOutput.Println(titles)
		if len(titles) <= 0 {
			continue
		}

		title := ""
		for _, t := range titles {
			if parseTime(t) == nil {
				title = t
			}
		}
		if title == "" {
			title = titles[0]
		}

		posts = append(posts, friends.FriendPost{
			Title: title,
			Link:  u.String(),
			Time:  toUnix(elementFindTime(item)),
		})
		if len(posts) >= 5 {
			break
		}
	}

	return posts
}
