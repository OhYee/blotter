package api

import (
	"github.com/OhYee/blotter/api/pkg/qiniu"
	"github.com/OhYee/blotter/register"
)

// GithubReposResponse response for GithubRepos api
type GetBucketsResponse struct {
	Buckets []string `json:"buckets"`
}

// GetBuckets get buckets name
func GetBuckets(context register.HandleContext) (err error) {

	res := new(GetBucketsResponse)
	if res.Buckets, err = qiniu.GetBuckets(); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

type GetImagesRequest struct {
	Bucket string `json:"bucket"`
	Prefix string `json:"prefix"`
	Marker string `json:"marker"`
	Number int    `json:"number"`
}

// GithubReposResponse response for GithubRepos api
type GetImagesResponse struct {
	Files   []*qiniu.File `json:"files"`
	Marker  string        `json:"marker"`
	HasNext bool          `json:"has_next"`
}

