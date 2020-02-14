package api

import (
	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/register"
)

// MarkdownRequest request of markdown api
type MarkdownRequest struct {
	Source string `json:"source"`
}

// MarkdownResponse response of markdown api
type MarkdownResponse struct {
	HTML string `json:"html"`
}

// Markdown render markdown to html
func Markdown(context *register.HandleContext) (err error) {
	args := new(MarkdownRequest)
	res := new(MarkdownResponse)
	context.RequestParams(args)

	if res.HTML, err = markdown.Render(args.Source); err != nil {
		return
	}

	context.ReturnJSON(res)
	return
}
