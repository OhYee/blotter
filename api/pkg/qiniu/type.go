package qiniu

import (
	"github.com/qiniu/api.v7/v7/storage"
)

type File struct {
	Name string `json:"name" bson:"name"`
	Size int64  `json:"size" bson:"size"`
	Time int64  `json:"time" bson:"time"`
}

func NewFileFromListItem(item storage.ListItem) *File {
	return &File{
		Name: "//static.oyohyee.com/" + item.Key,
		Size: item.Fsize,
		Time: item.PutTime,
	}
}
