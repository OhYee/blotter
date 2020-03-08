package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func rerender(filter bson.M) (err error) {
	posts := new([]map[string]string)
	_, err = mongo.Find("blotter", "posts", filter, options.Find().SetProjection(bson.M{
		"url": 1,
		"raw": 1,
	}), posts)
	if err != nil {
		return err
	}
	total := len(*posts)
	for idx, post := range *posts {
		fmt.Printf("%d/%d render %s\n", idx, total, post["url"])
		html, err := markdown.Render(post["raw"], true)
		if err != nil {
			return err
		}
		mongo.Update(
			"blotter", "posts",
			bson.M{
				"url": post["url"],
			},
			bson.M{
				"$set": bson.M{
					"content": html,
				},
			},
			nil,
		)
	}
	return
}

func main() {
	var err error
	var posts string
	flag.StringVar(&posts, "posts", "", "posts to rerender")
	flag.Parse()

	if posts == "*" {
		err = rerender(bson.M{})
	} else {
		urls := strings.Split(posts, ",")
		if len(urls) > 0 {
			err = rerender(bson.M{
				"url": bson.M{"$in": urls},
			})
		}
	}
	if err != nil {
		panic(err)
	}
}
