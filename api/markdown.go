package api

import (
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/goldmark-dot"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"bytes"
)

type MarkdownRequest struct {
	Source string `json:"source"`
}
type MarkdownResponse struct {
	HTML string `json:"html"`
}

func renderMarkdown(source string) (html string, err error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			dot.NewDot("dot-svg", highlighting.NewHTMLRenderer()),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(),
	)

	buf := bytes.NewBuffer([]byte{})
	if err = md.Convert([]byte(source), buf); err == nil {
		html = buf.String()
	}
	return
}

func Markdown(context *register.HandleContext) (err error) {
	args := new(MarkdownRequest)
	res := new(MarkdownResponse)
	context.RequestArgs(args)

	if res.HTML, err = renderMarkdown(args.Source); err != nil {
		return
	}

	context.ReturnJSON(res)
	return
}
