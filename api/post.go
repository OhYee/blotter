package api

import (
	"strings"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
)

// PostRequest request for post api
type PostRequest struct {
	URL string `json:"url"`
}

// Post get post by url
func Post(context *register.HandleContext) (err error) {
	args := PostRequest{}
	context.RequestParams(&args)

	res, err := post.GetPublicFieldPost(args.URL)
	if err != nil {
		return
	}

	if res.URL == args.URL && args.URL != "" {
		context.ReturnJSON(res)
	} else {
		context.PageNotFound()
	}
	return
}

// ViewRequest request for inc api
type ViewRequest struct {
	URL string `json:"url"`
}

// View view number of the url
func View(context *register.HandleContext) (err error) {
	args := PostRequest{}
	context.RequestParams(&args)

	go post.IncView(args.URL)
	return
}

// PostAdmin get posts with all fields
func PostAdmin(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := PostRequest{}
	context.RequestParams(&args)

	res, err := post.GetAllFieldPost(args.URL)
	if err != nil {
		return
	}

	go post.IncView(args.URL)

	if res.URL == args.URL && args.URL != "" {
		context.ReturnJSON(res)
	} else {
		context.PageNotFound()
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
	Search    string `json:"search"`
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
	context.RequestParams(args)

	res.Total, res.Posts, err = post.GetCardPosts(args.Offset, args.Number, args.Tag, args.SortField, args.SortType, args.Search)
	context.ReturnJSON(res)
	return
}

// PostsAdminRequest request of posts api
type PostsAdminRequest PostsRequest

// PostsAdminResponse response of posts api
type PostsAdminResponse struct {
	Total int64             `json:"total"`
	Posts []post.AdminField `json:"posts"`
}

// PostsAdmin get posts
func PostsAdmin(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := new(PostsAdminRequest)
	res := new(PostsAdminResponse)
	context.RequestParams(args)

	res.Total, res.Posts, err = post.GetAdminPosts(args.Offset, args.Number, args.Tag, args.SortField, args.SortType, args.Search)
	context.ReturnJSON(res)
	return
}

// PostExistedRequest request of PostExisted api
type PostExistedRequest struct {
	URL string `json:"url"`
}

// PostExistedResponse response of PostExistede api
type PostExistedResponse struct {
	Existed bool `json:"existed"`
}

// PostExisted return the post is existed
func PostExisted(context *register.HandleContext) (err error) {
	args := new(PostExistedRequest)
	res := new(PostExistedResponse)
	context.RequestParams(args)

	res.Existed = post.Existed(args.URL)

	context.ReturnJSON(res)
	return
}

type PostEditRequest struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	URL         string   `json:"url"`
	Abstract    string   `json:"abstract"`
	HeadImage   string   `json:"head_image"`
	View        int64    `json:"view"`
	PublishTime int64    `json:"publish_time"`
	EditTime    int64    `json:"edit_time"`
	Published   bool     `json:"published"`
	Raw         string   `json:"raw"`
	Tags        []string `json:"tags"`
}

func PostEdit(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := new(PostEditRequest)
	res := SimpleResponse{Success: true, Title: "操作成功"}

	context.RequestData(args)

	if args.ID == "" {
		err = post.NewPost(
			args.Title,
			args.Abstract,
			args.View,
			args.URL,
			args.PublishTime,
			args.EditTime,
			args.Raw,
			args.Tags,
			[]string{},
			args.Published,
			args.HeadImage,
		)
		if err != nil && strings.HasPrefix(err.Error(), "Post with url existed") {
			res.Success = false
			res.Title = "文章发布失败"
			res.Content = err.Error()
			err = nil
		}
	} else {
		err = post.UpdatePost(
			args.ID,
			args.Title,
			args.Abstract,
			args.View,
			args.URL,
			args.PublishTime,
			args.EditTime,
			args.Raw,
			args.Tags,
			[]string{},
			args.Published,
			args.HeadImage,
		)
	}
	if err != nil {
		return
	}

	context.ReturnJSON(res)
	return
}

// PostDeleteRequest request of PostDelete api
type PostDeleteRequest struct {
	ID string `json:"id"`
}

// PostDeleteResponse response of PostDeletee api
type PostDeleteResponse SimpleResponse

// PostDelete return the post is existed
func PostDelete(context *register.HandleContext) (err error) {
	if !user.CheckToken(context.GetCookie("token")) {
		context.Forbidden()
		return
	}

	args := new(PostDeleteRequest)
	res := PostDeleteResponse{
		Success: true,
		Title:   "删除成功",
	}
	context.RequestParams(args)

	post.Delete(args.ID)

	context.ReturnJSON(res)
	return
}
