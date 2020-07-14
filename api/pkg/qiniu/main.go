package qiniu

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

func getKeys() (accessKey, secretKey string) {
	var err error
	v, err := variable.Get("qiniu_access_key", "qiniu_secret_key")
	if err != nil {
		return
	}
	fmt.Println(v)
	if v.SetString("qiniu_access_key", &accessKey) != nil {
		return
	}
	if v.SetString("qiniu_secret_key", &secretKey) != nil {
		return
	}
	return
}

func GenerateToken() (token string) {
	accessKey, secretKey := getKeys()
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope:   "space",
		Expires: 60,
	}
	token = putPolicy.UploadToken(mac)
	return
}

func GetBuckets() (buckets []string, prefix []string, err error) {
	accessKey, secretKey := getKeys()

	var prefixString string
	v, err := variable.Get("qiniu_prefix")
	if err != nil {
		return
	}
	v.SetString("qiniu_prefix", &prefixString)
	prefix = strings.Split(prefixString, ",")

	mac := qbox.NewMac(accessKey, secretKey)
	bucketManager := storage.NewBucketManager(mac, &storage.Config{
		UseHTTPS: true,
	})
	buckets, err = bucketManager.Buckets(true)
	return
}

func GetImages(bucket string, prefix string, marker string, count int) (files []*File, next string, hasNext bool, err error) {
	accessKey, secretKey := getKeys()
	mac := qbox.NewMac(accessKey, secretKey)
	bucketManager := storage.NewBucketManager(mac, &storage.Config{
		UseHTTPS: true,
	})

	var staticDomain string
	v, err := variable.Get("qiniu_static_domain")
	if err != nil {
		return
	}
	v.SetString("qiniu_static_domain", &staticDomain)

	items := make([]storage.ListItem, 0)
	if items, _, next, hasNext, err = bucketManager.ListFiles(bucket, prefix, "", marker, count); err != nil {
		return
	}

	files = make([]*File, len(items))
	for i, item := range items {
		files[i] = NewFileFromListItem(item, staticDomain)
	}
	return
}

func DeleteImage(bucket string, key string) (err error) {
	accessKey, secretKey := getKeys()
	mac := qbox.NewMac(accessKey, secretKey)
	bucketManager := storage.NewBucketManager(mac, &storage.Config{
		UseHTTPS: true,
	})
	err = bucketManager.Delete(bucket, key)
	return
}
	return
}
