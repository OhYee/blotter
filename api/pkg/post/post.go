package post

import (
	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	pipeline := []bson.M{
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
	}
	if project != nil {
		pipeline = append(pipeline, bson.M{"$project": project})
	}
	return Query(
		pipeline,
		post,
	)
}

// GetAllFieldPost get all field post
func GetAllFieldPost(url string) (post CompleteField, err error) {
	postsDB := make([]CompleteFieldDB, 0)
	cnt, err := QueryByURL(url, &postsDB, nil)
	if cnt > 0 {
		post = postsDB[0].ToPost()
	}
	return
}

// GetPublicFieldPost get all field post
func GetPublicFieldPost(url string) (post PublicField, err error) {
	postsDB := make([]PublicFieldDB, 0)
	cnt, err := QueryByURL(url, &postsDB, nil)
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

func Existed(url string) bool {
	cnt, err := mongo.Find("blotter", "posts", bson.M{
		"url": url,
	}, nil, nil)
	return !(err == nil && cnt == 0)
}

// NewPost insert a new post to database
func NewPost(
	title string,
	abstract string,
	view int64,
	url string,
	publishTime int64,
	editTime int64,
	raw string,
	tags []string,
	keywords []string,
	published bool,
	headImage string,
) (err error) {

	html, err := markdown.Render(raw)
	if err != nil {
		return
	}
	tagsID := make([]primitive.ObjectID, 0)
	for idx, tagName := range tags {
		if tagsID[idx], err = primitive.ObjectIDFromHex(tagName); err != nil {
			return
		}
	}

	p := DB{}
	p.ID = primitive.NewObjectID()
	p.Title = title
	p.Abstract = abstract
	p.View = view
	p.URL = url
	p.PublishTime = publishTime
	p.EditTime = editTime
	p.Content = html
	p.Raw = raw
	p.Tags = tagsID
	p.Keywords = keywords
	p.Published = published
	p.HeadImage = headImage

	if Existed(url) {
		err = errors.New("Post with url existed: %s", url)
		return
	}

	_, err = mongo.Add("blotter", "posts", nil, p)
	return
}

// UpdatePost update post data of id
func UpdatePost(
	id string,
	title string,
	abstract string,
	view int64,
	url string,
	publishTime int64,
	editTime int64,
	raw string,
	tags []string,
	keywords []string,
	published bool,
	headImage string,
) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	html, err := markdown.Render(raw)
	if err != nil {
		return
	}
	tagsID := make([]primitive.ObjectID, 0)
	for idx, tagName := range tags {
		if tagsID[idx], err = primitive.ObjectIDFromHex(tagName); err != nil {
			return
		}
	}

	p := DB{}
	p.ID = primitive.NewObjectID()
	p.Title = title
	p.Abstract = abstract
	p.View = view
	p.URL = url
	p.PublishTime = publishTime
	p.EditTime = editTime
	p.Content = html
	p.Raw = raw
	p.Tags = tagsID
	p.Keywords = keywords
	p.Published = published
	p.HeadImage = headImage
	_, err = mongo.Update("blotter", "posts", bson.M{
		"_id": objectID,
	}, p, nil)
	return
}
