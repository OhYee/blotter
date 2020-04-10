package api

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/api/pkg/post"
	"github.com/OhYee/blotter/api/pkg/tag"
	"github.com/OhYee/blotter/api/pkg/user"
	"github.com/OhYee/blotter/register"
)

// TagsRequest request for api tags
type TagsRequest struct {
	Keyword   string `json:"keyword"`
	Offset    int64  `json:"offset"`
	Number    int64  `json:"number"`
	SortField string `json:"sort_field"`
	SortInc   bool   `json:"sort_inc"`
}

// TagsResponse request without keyword
type TagsResponse struct {
	Total int64           `json:"total"`
	Tags  []tag.WithCount `json:"tags"`
}

// Tags query tags
func Tags(context register.HandleContext) (err error) {
	args := new(TagsRequest)
	res := new(TagsResponse)

	context.RequestArgs(args)

	res.Total, res.Tags, err = tag.GetTags(args.Keyword, args.Offset, args.Number, args.SortField, args.SortInc)
	if err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

// TagEditRequest request of TagEdit api
type TagEditRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Short string `json:"short"`
	Color string `json:"color"`
	Icon  string `json:"icon"`
}

// TagEditResponse response of TagEdit api
type TagEditResponse SimpleResponse

// TagEdit edit tag info
func TagEdit(context register.HandleContext) (err error) {
	if !user.CheckUserPermission(context) {
		context.Forbidden()
		return
	}

	args := new(TagEditRequest)
	res := new(TagEditResponse)
	context.RequestArgs(args)

	res.Success = true
	if args.ID == "" {
		res.Title = ""
		if err = tag.New(args.Name, args.Short, args.Color, args.Icon); err == nil {
			res.Title = "新标签添加成功"
		} else {
			res.Success = false
			res.Title = "新标签添加失败"
			res.Content = err.Error()
		}

	} else {
		if err = tag.Update(args.ID, args.Name, args.Short, args.Color, args.Icon); err == nil {
			res.Title = "修改标签成功"
		} else {
			res.Success = false
			res.Title = "修改标签失败"
			res.Content = err.Error()
		}
	}

	err = context.ReturnJSON(res)
	return
}

// TagDeleteRequest request of TagDelete api
type TagDeleteRequest struct {
	ID string `json:"id"`
}

// TagDeleteResponse response of TagDelete api
type TagDeleteResponse SimpleResponse

// TagDelete edit tag info
func TagDelete(context register.HandleContext) (err error) {
	if !user.CheckUserPermission(context) {
		context.Forbidden()
		return
	}

	args := new(TagDeleteRequest)
	res := new(TagDeleteResponse)
	context.RequestArgs(args)

	if err = tag.Delete(args.ID); err != nil {
		res.Success = false
		res.Title = "标签删除失败"
		res.Content = err.Error()
	} else {
		res.Success = true
		res.Title = "标签删除成功"
	}

	err = context.ReturnJSON(res)
	return
}

// TagExistedRequest request of TagExisted api
type TagExistedRequest struct {
	ID    string `json:"id"`
	Short string `json:"short"`
}

// TagExistedResponse response of TagExisted api
type TagExistedResponse struct {
	Existed bool `json:"existed"`
}

// TagExisted edit tag info
func TagExisted(context register.HandleContext) (err error) {
	args := new(TagExistedRequest)
	res := new(TagExistedResponse)
	context.RequestArgs(args)

	if res.Existed, err = tag.Existed(args.ID, args.Short); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

// TagRequest request of Tag api
type TagRequest struct {
	Number int64  `json:"number"`
	Offset int64  `json:"offset"`
	Tag    string `json:"tag"`
}

// TagResponse response of Tag api
type TagResponse struct {
	Tag   tag.Type         `json:"tag"`
	Total int64            `json:"total"`
	Posts []post.CardField `json:"posts"`
}

// Tag edit tag info
func Tag(context register.HandleContext) (err error) {
	args := new(TagRequest)
	res := new(TagResponse)
	context.RequestArgs(args)

	if res.Tag, err = tag.Get(args.Tag); err != nil {
		if strings.HasPrefix(err.Error(), fmt.Sprintf("No tag %s", args.Tag)) {
			context.PageNotFound()
		}
		return
	}
	if res.Total, res.Posts, err = post.GetCardPosts(
		args.Offset, args.Number,
		[]string{res.Tag.ID}, []string{},
		"", -1, "", []string{},
	); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}
