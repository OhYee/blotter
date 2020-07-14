package qiniu

import (
	"github.com/qiniu/api.v7/v7/storage"
)

type File struct {
	Key  string `json:"key" bson:"key"`
	Link string `json:"link" bson:"link"`
	Size int64  `json:"size" bson:"size"`
	Time int64  `json:"time" bson:"time"`
}

func NewFileFromListItem(item storage.ListItem, domain string) *File {

	return &File{
		Key:  item.Key,
		Link: storage.MakePublicURL(domain, item.Key),
		Size: item.Fsize,
		Time: item.PutTime,
	}
}
