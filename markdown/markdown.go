package main

import (
	"github.com/OhYee/goldmark-dot"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"bytes"
)

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
