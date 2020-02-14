package api

import (
	"github.com/OhYee/blotter/api/pkg/tag"
	"github.com/OhYee/blotter/register"
)

// TagsRequest request for api tags
type TagsRequest struct {
	Keyword string `json:"keyword"`
}

// TagsResponse request without keyword
type TagsResponse struct {
	Total int64      `json:"total"`
	Tags  []tag.Type `json:"tags"`
}

// TagsResponseWithCount request with keyword
type TagsResponseWithCount struct {
	Total int64           `json:"total"`
	Tags  []tag.WithCount `json:"tags"`
}

// Tags query tags
func Tags(context *register.HandleContext) (err error) {
	args := new(TagsRequest)
	var res interface{}

	context.RequestParams(args)

	if len(args.Keyword) == 0 {
		resWithCount := new(TagsResponseWithCount)
		resWithCount.Total, resWithCount.Tags, err = tag.GetTags()
		res = resWithCount
	} else {
		resWithoutCount := new(TagsResponse)
		resWithoutCount.Total, resWithoutCount.Tags, err = tag.SearchTags(args.Keyword)
		res = resWithoutCount
	}
	if err != nil {
		return
	}
	err = context.ReturnJSON(res)
	if err != nil {
		return
	}
	return
}
