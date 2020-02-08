package api

import (
	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/register"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

// PostRequest request for post api
type PostRequest struct {
	URL string `json:"url"`
}

// Post get post by url
func Post(context *register.HandleContext) (err error) {
	args := PostRequest{}
	context.RequestArgs(&args)

	res, err := post.GetPublicFieldPost(args.URL)
	if err != nil {
		return
	}

	go post.IncView(args.URL)

	if res.URL != args.URL {
		context.ReturnJSON(res)
	} else {
		context.Response.WriteHeader(404)
	}
	return
}

// PostsRequest request of posts api
type PostsRequest struct {
	Number    int64  `json:"number"`
	Offset    int64  `json:"offset"`
	Tag       string `json:"tag"`
	SortField string `json:"sort_field"`
	SortType  int    `json:"sort_type"`
}

// PostsResponse response of posts api
type PostsResponse struct {
	Total int64            `json:"total"`
	Posts []post.CardField `json:"posts"`
}

// Posts get posts
func Posts(context *register.HandleContext) (err error) {
	args := new(PostsRequest)
	res := new(PostsResponse)
	context.RequestArgs(args)

	res.Total, res.Posts, err = post.GetCardPosts(args.Offset, args.Number, args.Tag, args.SortField, args.SortType)
	context.ReturnJSON(res)
	return
}
