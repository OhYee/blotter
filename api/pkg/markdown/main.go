package markdown

import (
	dot "github.com/OhYee/goldmark-dot"
	qjskatex "github.com/graemephi/goldmark-qjs-katex"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"

	"bytes"
)

func Render(source string, renderHTML bool) (htmlResult string, err error) {
	renderOpts := []renderer.Option{html.WithHardWraps()}
	if renderHTML {
		renderOpts = append(renderOpts, html.WithUnsafe())
	}

	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.DefinitionList,
			extension.Footnote,
			dot.NewDot("dot-svg", highlighting.NewHTMLRenderer()),
			// mathjax.MathJax,
			&qjskatex.Extension{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(renderOpts...),
	)

	buf := bytes.NewBuffer([]byte{})
	if err = md.Convert([]byte(source), buf); err == nil {
		htmlResult = buf.String()
	}
	return
}
