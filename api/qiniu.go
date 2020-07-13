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
