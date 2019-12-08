package main

import (
	"context"

	"github.com/OhYee/blotter/markdown/proto"
	"github.com/micro/go-micro"
)


func renderMarkdown(source string) (string, error) {
	service := micro.NewService(micro.Name("site.client"))
	service.Init()
	markdown := proto.NewMarkdownService("markdown", service.Client())
	rsp, err := markdown.RenderMarkdown(context.TODO(), &proto.RenderMarkdownRequest{Source: source})
	if err != nil {
		return "", err
	}
	return rsp.Html, nil
}