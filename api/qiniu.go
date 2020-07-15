package api

import (
	"github.com/OhYee/blotter/api/pkg/qiniu"
	"github.com/OhYee/blotter/register"
	"github.com/OhYee/rainbow/errors"
)

// GithubReposResponse response for GithubRepos api
type GetBucketsResponse struct {
	Buckets []string `json:"buckets"`
	Prefix  []string `json:"prefix"`
}

// GetBuckets get buckets name
func GetBuckets(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	res := new(GetBucketsResponse)
	if res.Buckets, res.Prefix, err = qiniu.GetBuckets(); err != nil {
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
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

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
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	res := new(GetQiniuTokenResponse)

	res.Token = qiniu.GenerateToken()

	err = context.ReturnJSON(res)
	return
}

type DeleteImageRequest struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
}
type DeleteImageResponse SimpleResponse

// DeleteImage get images of bucket
func DeleteImage(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(DeleteImageRequest)
	res := new(DeleteImageResponse)
	context.RequestArgs(args)

	if err = qiniu.DeleteImage(args.Bucket, args.Key); err != nil {
		res.Success = false
		res.Title = "删除失败"
		res.Content = errors.ShowStack(err)
	} else {
		res.Success = true
		res.Title = "删除成功"
	}

	err = context.ReturnJSON(res)
	return
}

type RenameImageRequest struct {
	Bucket string `json:"bucket"`
	Key    string `json:"key"`
	NewKey string `json:"new_key"`
}
type RenameImageResponse SimpleResponse

// RenameImage get images of bucket
func RenameImage(context register.HandleContext) (err error) {
	if !context.GetUser().HasPermission() {
		context.Forbidden()
		return
	}

	args := new(RenameImageRequest)
	res := new(RenameImageResponse)
	context.RequestArgs(args)

	if err = qiniu.RenameImage(args.Bucket, args.Key, args.NewKey); err != nil {
		res.Success = false
		res.Title = "重命名失败"
		res.Content = errors.ShowStack(err)
	} else {
		res.Success = true
		res.Title = "重命名成功"
	}

	err = context.ReturnJSON(res)
	return
}
