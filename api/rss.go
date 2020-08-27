package api

import (
	"fmt"
	"strings"
	"time"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"

	gt "github.com/OhYee/goutils/time"
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

	total, posts, err := post.GetCardPosts(0, 0, []string{}, []string{}, "", 0, "", []string{})
	if err != nil {
		return
	}

	data := make([]string, total)
	for idx, post := range posts {
		datetime := post.PublishTime
		if t, e := time.ParseInLocation("2006-01-02 15:04:05", post.PublishTime, gt.ChinaTimeZone); e == nil {
			datetime = t.Format("Mon, 02 Jan 2006 15:04:05 -0700")
		}

		data[idx] = fmt.Sprintf(
			postFormat,
			post.Title,
			fmt.Sprintf("%s/post/%s", root, post.URL),
			post.Abstract,
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
