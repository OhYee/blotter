package spider

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/OhYee/blotter/api/pkg/friends"
	"github.com/OhYee/blotter/output"
)

func toInt64(v interface{}) int64 {
	switch t := v.(type) {
	case uint, uint8, uint16, uint32, uint64:
		return int64(reflect.ValueOf(t).Uint())
	case int, int8, int16, int32, int64:
		return reflect.ValueOf(t).Int()
	case float32, float64:
		return int64(reflect.ValueOf(t).Float())
	case string:
		if i, err := strconv.ParseInt(t, 10, 64); err == nil {
			return i
		}
	}
	return 0
}

type jsonObject = map[string]interface{}
type jsonArray = []interface{}

func rangeJSON(data interface{}) []friends.FriendPost {
	posts := make([]friends.FriendPost, 0)
	switch v := data.(type) {
	case jsonObject:
		posts = append(posts, rangeObject(v)...)
	case jsonArray:
		posts = append(posts, rangeArray(v)...)
	}
	return posts
}

func rangeObject(data jsonObject) []friends.FriendPost {
	posts := make([]friends.FriendPost, 0)
	if _, ok := data["title"]; ok {
		title, link, ts := parsePostObject(data)
		if len(title) > 0 && len(link) > 0 {
			timestamp := int64(0)
			if ts != nil {
				timestamp = ts.Unix()
			}
			posts = append(posts, friends.FriendPost{
				Title: title,
				Link:  link,
				Time:  timestamp,
			})
		}

	}
	for _, v := range data {
		posts = append(posts, rangeJSON(v)...)
	}
	return posts
}

func rangeArray(data []interface{}) []friends.FriendPost {
	posts := make([]friends.FriendPost, 0)
	for _, item := range data {
		posts = append(posts, rangeJSON(item)...)
	}
	return posts
}

func parsePostObject(data jsonObject) (title, link string, ts *time.Time) {
	switch v := data["title"].(type) {
	case string:
		title = v
	case jsonObject:
		for _, value := range v {
			if valueString, isString := value.(string); isString {
				title += valueString
			}
		}
	}

	if v, ok := data["title"]; ok {
		if _, ok2 := v.(string); ok2 {
			title = v.(string)
		}
	}
	for _, linkKey := range linkKeys {
		if linkValue, ok := data[linkKey]; ok {
			if v, ok2 := linkValue.(string); ok2 {
				link = v
				break
			}
		}
	}

	for _, timeKey := range timeKeys {
		if timeValue, ok := data[timeKey]; ok {
			ts = parseTime(timeValue)
			if ts != nil {
				break
			}
		}
	}

	if ts == nil {
		for _, value := range data {
			ts = parseTime(value)
			if ts != nil {
				break
			}
		}
	}

	return
}

func readJSON(link, content string) (posts []friends.FriendPost) {
	output.DebugOutput.Println(link, "readJSON")
	posts = make([]friends.FriendPost, 0)

	if len(content) > 0 && content[0] == '[' {
		content = fmt.Sprintf("{\"data\":%s}", content)
	}
	data := make(jsonObject)
	if err := json.Unmarshal([]byte(content), &data); err != nil {
		output.ErrOutput.Println(err)
		return
	}

	if _, ok := data["posts"]; ok {
		posts = append(posts, rangeJSON(data["posts"])...)
	} else {
		posts = append(posts, rangeJSON(data)...)
	}

	return
}
