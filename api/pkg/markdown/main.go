package markdown

import (
	"fmt"

	"github.com/OhYee/blotter/output"
	dot "github.com/OhYee/goldmark-dot"
	ext "github.com/OhYee/goldmark-fenced_codeblock_extension"
	img "github.com/OhYee/goldmark-image"
	uml "github.com/OhYee/goldmark-plantuml"
	python "github.com/OhYee/goldmark-python"
	qjskatex "github.com/graemephi/goldmark-qjs-katex"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer"
	"github.com/yuin/goldmark/renderer/html"

	"bytes"
)

var exts = ext.NewExt(
	ext.RenderMap{
		Languages:      []string{"dot-svg"},
		RenderFunction: dot.NewDot(50, "dot-svg").Renderer,
	},
	ext.RenderMap{
		Languages:      []string{"uml-svg"},
		RenderFunction: uml.NewUML(50, "uml-svg").Renderer,
	},
	ext.RenderMap{
		Languages:      []string{"python-output"},
		RenderFunction: python.NewPython(50, "python3", "python-output").Renderer,
	},
	ext.RenderMap{
		Languages: []string{"*"},
		RenderFunction: ext.GetFencedCodeBlockRendererFunc(
			highlighting.NewHTMLRenderer(
				highlighting.WithGuessLanguage(true),
				highlighting.WithStyle("trac"),
			),
		),
	},
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
			exts,
			img.NewImg("image", func(args img.ImgArgs, class string, renderImg img.RenderImgFunc) string {
				var title = ""
				if args.Title != "" {
					title = fmt.Sprintf("<span class=\"img-title\">%s</span>", args.Title)
				} else if args.Alt != "" {
					title = fmt.Sprintf("<span class=\"img-title\">%s</span>", args.Alt)
				}
				return fmt.Sprintf(
					"<div class='%s'><a href='%s' target='_blank' rel='noopener noreferrer'>%s</a>%s</div>",
					class, args.Src, renderImg(args), title,
				)
			}),
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
