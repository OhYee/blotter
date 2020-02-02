package main

import (
	"flag"
	"github.com/OhYee/blotter/api"
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
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
	for _, post := range *posts {
		html, err := api.RenderMarkdown(post["raw"])
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
	flag.StringVar(&posts, "posts", "*", "posts to rerender")
	flag.Parse()

	if posts == "*" {
		err = rerender(bson.M{})
	} else {
		err = rerender(bson.M{
			"url": bson.M{"$in": strings.Split(posts, ",")},
		})
	}
	if err != nil {
		panic(err)
	}
}
