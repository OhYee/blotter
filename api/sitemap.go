package api

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/register"
)

// SitemapTXT sitemap.txt
func SitemapTXT(context *register.HandleContext) (err error) {
	variables, err := variable.Get("root")
	if err != nil {
		return
	}
	root, _ := variables.GetString("root")
	total, posts, err := post.GetCardPosts(0, 0, []string{}, []string{}, "", 0, "", []string{})
	if err != nil {
		return
	}

	data := make([]string, total)
	for idx, post := range posts {
		data[idx] = fmt.Sprintf("%s/post/%s", root, post.URL)
	}

	context.ReturnText(strings.Join(data, "\n"))
	return
}

// SitemapXML sitemap.txt
func SitemapXML(context *register.HandleContext) (err error) {
	variables, err := variable.Get("root")
	if err != nil {
		return
	}
	root, _ := variables.GetString("root")
	total, posts, err := post.GetCardPosts(0, 0, []string{}, []string{}, "", 0, "", []string{})
	if err != nil {
		return
	}

	data := make([]string, total)
	for idx, post := range posts {
		data[idx] = fmt.Sprintf("<url><loc>%s/post/%s</loc></url>", root, post.URL)
	}

	context.ReturnXML(
		fmt.Sprintf(
			"<?xml version=\"1.0\"?><urlset xmlns=\"http://www.sitemaps.org/schemas/sitemap/0.9\">%s</urlset>",
			strings.Join(data, "\n"),
		),
	)
	return
}
