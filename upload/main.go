package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
)

func main() {
	accessKey := "cNqMXY9XLBDc6XD5D91GV4M6NMb-fmn1_dEjtReS"
	secretKey := "fDcsKn3dfTXKMs4OfF_OZUxbZ-uHbrHEbMyc9hXj"
	mac := qbox.NewMac(accessKey, secretKey)
	// cfg := storage.Config{
	// 	// 是否使用https域名进行资源管理
	// 	UseHTTPS: false,
	// }
	// 指定空间所在的区域，如果不指定将自动探测
	// 如果没有特殊需求，默认不需要指定
	//cfg.Zone=&storage.ZoneHuabei
	// bucketManager := storage.NewBucketManager(mac, &cfg)
	// cfg.

	// a, _ := bucketManager.GetBucketInfo("space")
	// fmt.Printf("%+v\n", a)

	// fmt.Println(bucketManager.Buckets(true))
	// a, b, c, d, e := bucketManager.ListFiles("space", "", "", "", 1000)
	// fmt.Printf("%+v\n", a)
	// fmt.Printf("%+v\n", b)
	// fmt.Printf("%+v\n", c)
	// fmt.Printf("%+v\n", d)
	// fmt.Printf("%+v\n", e)

	// for aa := range a {
	// 	fmt.Println(a[aa])
	// }

	putPolicy := storage.PutPolicy{
		Scope: "space",
	}
	putPolicy.Expires = 7200 //示例2小时有效期

	upToken := putPolicy.UploadToken(mac)
	fmt.Println(upToken)

	client := &http.Client{}
	reqest, _ := http.NewRequest("GET", "https://api.qiniu.com/v6/domain/list?tbl=space", nil)
	reqest.Header.Add("Authorization", fmt.Sprintf("Qiniu %s", upToken))
	response, _ := client.Do(reqest)
	bb, _ := ioutil.ReadAll(response.Body)
	fmt.Printf("%s\n", bb)
}
