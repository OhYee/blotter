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

// GetImages get images of bucket
func GetImages(context register.HandleContext) (err error) {
	args := new(GetImagesRequest)
	res := new(GetImagesResponse)
	context.RequestArgs(args)

	if res.Files, res.Marker, res.HasNext, err = qiniu.GetImages(args.Bucket, args.Prefix, args.Marker, args.Number); err != nil {
		return
	}

	err = context.ReturnJSON(res)
	return
}

type GetQiniuTokenResponse struct {
	Token string `json:"token"`
}

// GetQiniuToken get images of bucket
func GetQiniuToken(context register.HandleContext) (err error) {
	res := new(GetQiniuTokenResponse)

	res.Token = qiniu.GenerateToken()

	err = context.ReturnJSON(res)
	return
}
