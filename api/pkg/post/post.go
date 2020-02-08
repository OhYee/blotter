package post

import (
	"github.com/OhYee/blotter/mongo"
	"go.mongodb.org/mongo-driver/bson"
)

// Query posts of url
func Query(pipeline []bson.M, res interface{}) (total int64, err error) {
	total, err = mongo.Aggregate(
		"blotter", "posts",
		pipeline,
		nil, res,
	)
	if err != nil {
		return
	}
	return
}

// QueryByURL query post by url
func QueryByURL(url string, post interface{}, project bson.M) (total int64, err error) {
	return Query(
		[]bson.M{
			{
				"$match": bson.M{
					"url": url,
				},
			},
			{
				"$limit": 1,
			},
			{
				"$lookup": bson.M{
					"from":         "tags",
					"localField":   "tags",
					"foreignField": "_id",
					"as":           "tags",
				},
			},
			{
				"$project": project,
			},
		},
		post,
	)
}

// GetAllFieldPost get all field post
func GetAllFieldPost(url string) (post CompleteField, err error) {
	postsDB := make([]CompleteFieldDB, 0)
	cnt, err := QueryByURL(url, &postsDB, bson.M{})
	if cnt > 0 {
		post = postsDB[0].ToPost()
	}
	return
}

// GetPublicFieldPost get all field post
func GetPublicFieldPost(url string) (post PublicField, err error) {
	postsDB := make([]PublicFieldDB, 0)
	cnt, err := QueryByURL(url, &postsDB, bson.M{})
	if cnt > 0 {
		post = postsDB[0].ToPost()
	}
	return
}

// IncView view +1
func IncView(url string) {
	mongo.Update(
		"blotter", "posts", bson.M{"url": url},
		bson.M{"$inc": bson.M{"view": 1}}, nil,
	)
}

// GetCardPosts get all field post
func GetCardPosts(offset int64, number int64, tag string, sortField string, sortType int) (total int64, posts []CardField, err error) {
	sort := bson.M{"$sort": bson.M{"publish_time": -1}}
	if sortField != "" && (sortType == 1 || sortType == -1) {
		sort = bson.M{"$sort": bson.M{sortField: sortType}}
	}
	pipeline := []bson.M{
		sort,
		{
			"$lookup": bson.M{
				"from":         "tags",
				"localField":   "tags",
				"foreignField": "_id",
				"as":           "tags",
			},
		},
	}

	if tag != "" {
		pipeline = append(
			pipeline,
			bson.M{"$set": bson.M{"temp": "$tags.short"}},
			bson.M{"$match": bson.M{"temp": tag}},
		)
	}

	if offset != 0 || number != 0 {
		pipeline = append(
			pipeline,
			bson.M{"$limit": offset + number},
			bson.M{"$skip": offset},
		)
	}

	postsDB := make([]CardFieldDB, 0)
	total, err = Query(pipeline, &postsDB)
	posts = make([]CardField, len(postsDB))
	for idx, post := range postsDB {
		posts[idx] = post.ToCard()
	}
	return
}

// func NewPostDatabase(title string, abstract string, url string, raw string, tags []string, keywords []string, published bool, headImage string) *PostDatabase {
// 	html, err := RenderMarkdown(raw)
// 	if err != nil {
// 		html = raw
// 	}
// 	ids := make([]struct {
// 		ID primitive.ObjectID `bson:"_id"`
// 	}, 0)

// 	mongo.Aggregate("blotter", "tags", []bson.M{
// 		{
// 			"$match": bson.M{
// 				"short": bson.M{"$in": tags},
// 			},
// 		},
// 		{
// 			"$project": bson.M{"_id": 1},
// 		},
// 	}, nil, &ids)

// 	tagIDs := make([]primitive.ObjectID, len(ids))
// 	for idx, tag := range ids {
// 		tagIDs[idx] = tag.ID
// 	}

// 	return &PostDatabase{
// 		ID:          primitive.NewObjectID(),
// 		Title:       title,
// 		Abstract:    abstract,
// 		View:        0,
// 		URL:         url,
// 		Raw:         raw,
// 		PublishTime: time.Now().Unix(),
// 		EditTime:    0,
// 		Content:     html,
// 		Tags:        tagIDs,
// 		Keywords:    keywords,
// 		Published:   published,
// 		HeadImage:   headImage,
// 	}
// }
