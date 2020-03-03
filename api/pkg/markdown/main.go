package markdown

import (
	"github.com/OhYee/blotter/output"
	dot "github.com/OhYee/goldmark-dot"
	ext "github.com/OhYee/goldmark-fenced_codeblock_extension"
	img "github.com/OhYee/goldmark-image"
	uml "github.com/OhYee/goldmark-plantuml"
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
			ext.NewExt(
				ext.RenderMap{
					Language:       []string{"dot-svg"},
					RenderFunction: dot.NewDot("dot-svg").Renderer,
				},
				ext.RenderMap{
					Language:       []string{"uml-svg"},
					RenderFunction: uml.NewUML("uml-svg").Renderer,
				},
				ext.RenderMap{
					Language:       []string{"*"},
					RenderFunction: ext.GetFencedCodeBlockRendererFunc(highlighting.NewHTMLRenderer()),
				},
			),
			img.NewImg("image", nil),
			&qjskatex.Extension{},
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
			parser.WithAttribute(),
		),
		goldmark.WithRendererOptions(renderOpts...),
	)

	buf := bytes.NewBuffer([]byte{})
	if err = md.Convert([]byte(source), buf); err == nil {
		htmlResult = buf.String()
	} else {
		output.Err(err)
	}
	return
}
