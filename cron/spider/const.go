package spider

import "time"

const (
	UserAgent = "OhYee-Spider"
	Timeout   = 120 * time.Second

	year2000 = 946656000 // 2000-01-01 00:00:00
)

var (
	linkKeys = []string{"link", "href", "url"}
	timeKeys = []string{
		"date", "time",
		"pub", "pubdate", "pub_date", "pubtime", "pub_time",
		"publish", "publishdate", "publish_date", "publish_time", "publish_time",
		"published", "publisheddate", "published_date", "publishedtime", "published_time",
		"create", "createdate", "create_date", "createtime", "create_time",
		"created", "createddate", "created_date", "createdtime", "created_time",
		"update", "updatedate", "update_date", "updatetime", "update_time",
		"updated", "updateddate", "updated_date", "updatedtime", "updated_time",
		"modify", "modifydate", "modify_date", "modifytime", "modify_time",
		"modified", "modifieddate", "modified_date", "modifiedtime", "modified_time",
	}
)
