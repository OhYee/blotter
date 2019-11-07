package main

import (
	"context"
	"fmt"

	"github.com/OhYee/blotter/markdown/proto"
	"github.com/OhYee/blotter/markdown/markdown"
	"github.com/micro/go-micro"
)

//go:generate protoc --proto_path=$GOPATH/src:. --micro_out=. --go_out=. ./proto/proto.proto

// Markdown server object
type Markdown struct{}

// RenderMarkdown from markdown to html
func (g *Markdown) RenderMarkdown(ctx context.Context, req *proto.RenderMarkdownRequest, rsp *proto.RenderMarkdownResponse) (err error) {
	rsp.Html, err = markdown.RenderMarkdown(req.Source)
	return
}

func main() {
	service := micro.NewService(
		micro.Name("markdown"),
	)
	service.Init()
	proto.RegisterMarkdownHandler(service.Server(), new(Markdown))
	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
