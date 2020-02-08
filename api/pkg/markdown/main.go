package markdown

import (
	dot "github.com/OhYee/goldmark-dot"
	qjskatex "github.com/graemephi/goldmark-qjs-katex"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"

	"bytes"
)

func Render(source string) (html string, err error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			dot.NewDot("dot-svg", highlighting.NewHTMLRenderer()),
			// mathjax.MathJax,
			&qjskatex.Extension{},
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
