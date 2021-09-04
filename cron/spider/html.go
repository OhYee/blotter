package spider

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/output"
	"github.com/chromedp/chromedp"
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

func getChromePath() string {
	env := os.Environ()
	for _, s := range env {
		ss := strings.Split(s, "=")
		if len(ss) >= 2 {
			key := ss[0]
			value := strings.Join(ss[1:], "")
			if strings.ToUpper(key) == "CHROME_PATH" {
				return value
			}
		}
	}
	return ""
}
func GetHTML(u string) string {
	opts := []func(*chromedp.ExecAllocator){
		chromedp.Flag("headless", true),
		chromedp.Flag("blink-settings", "imagesEnabled=false"),
		chromedp.UserAgent(UserAgent),
		chromedp.Flag("ignore-certificate-errors", true),
	}

	chromePath := getChromePath()
	if len(chromePath) > 0 {
		opts = append(
			opts,
			chromedp.ExecPath(
				fmt.Sprintf(
					"%s/%s",
					strings.TrimRight(chromePath, "/"),
					"chrome",
				),
			),
		)
	}

	ctx, cancel := chromedp.NewExecAllocator(
		context.Background(),

		opts...,
	)
	defer cancel()

	ctx2, cancel2 := chromedp.NewContext(
		ctx,
	)
	defer cancel2()

	ctx3, cancel3 := context.WithTimeout(ctx2, Timeout)
	defer cancel3()

	var res string
	err := chromedp.Run(
		ctx3,
		chromedp.Navigate(u),
		chromedp.OuterHTML("html", &res, chromedp.ByQuery),
	)
	if err != nil {
		output.ErrOutput.Println(u, err)
	}
	return res
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

func ReadHTML(u string) []friends.FriendPost {
	hostURL, _ := url.Parse(u)

	c := GetHTML(u)
	doc, _ := html.Parse(bytes.NewBufferString(c))
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
		if len(titles) <= 0 {
			continue
		}
		posts = append(posts, friends.FriendPost{
			Title: titles[0],
			Link:  u.String(),
			Time:  toUnix(elementFindTime(item)),
		})
		if len(posts) >= 5 {
			break
		}
	}

	return posts
}
