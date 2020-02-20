package api

import (
	"fmt"
	"strings"

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
func RSSXML(context *register.HandleContext) (err error) {
	variables, err := variable.Get("root", "email", "blog_name", "author")
	if err != nil {
		return
	}
	root, _ := variables.GetString("root")
	email, _ := variables.GetString("email")
	blogName, _ := variables.GetString("blog_name")
	author, _ := variables.GetString("author")

	total, posts, err := post.GetCardPosts(0, 0, "", "", 0, "")
	if err != nil {
		return
	}

	data := make([]string, total)
	for idx, post := range posts {
		data[idx] = fmt.Sprintf(
			postFormat,
			post.Title,
			fmt.Sprintf("%s/post/%s", root, post.URL),
			post.Abstract,
			post.PublishTime,
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