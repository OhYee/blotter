package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"
)

const (
	rssFormat = `<?xml version="1.0" ?>
<rss version="2.0">
    <channel>
        <title>%s</title>
        <link>%s</link>
        <description>%s</description>
        <language>zh-cn</language>
        <webMaster>%s</webMaster>
        <image>
            <url>%s</url>
            <title>%s</title>
            <link>%s</link>
        </image>
		%s
    </channel>
</rss>`
	postFormat = `<item>
            <title>%s</title>
            <link>%s</link>
            <description>%s</description>
            <pubDate>%s</pubDate>
        </item>`
)

// RSSXML sitemap.txt
func RSSXML(context register.HandleContext) (err error) {
	variables, err := variable.Get("root", "email", "blog_name", "author")
	if err != nil {
		return
	}
	root, _ := variables.GetString("root")
	email, _ := variables.GetString("email")
	blogName, _ := variables.GetString("blog_name")
	author, _ := variables.GetString("author")
	withoutTags := []string{}

	total, posts, err := post.GetCardPosts(0, 0, []string{}, withoutTags, "", 0, "", []string{}, true)

	if err != nil {
		return
	}

	data := make([]string, total)
	for idx, post := range posts {
		t := time.Unix(post.PublishTime, 0).Local()
		datetime := t.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		data[idx] = fmt.Sprintf(
			postFormat,
			htmlEscape(post.Title),
			fmt.Sprintf("%s/post/%s", root, post.URL),
			htmlEscape(post.Abstract),
			datetime,
		)
	}

	context.ReturnXML(
		fmt.Sprintf(
			rssFormat,
			blogName,
			root,
			fmt.Sprintf("%s by %s", blogName, author),
			email,
			fmt.Sprintf("%s/static/img/logo.svg", root),
			blogName,
			root,
			strings.Join(data, "\n"),
		),
	)
	return
}

func htmlEscape(s string) string {
	s = strings.Replace(s, "&", "&amp;", -1)
	s = strings.Replace(s, "<", "&lt;", -1)
	s = strings.Replace(s, ">", "&gt;", -1)
	return s
}
