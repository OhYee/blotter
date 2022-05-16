package post

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/OhYee/blotter/api/pkg/markdown"
	"github.com/OhYee/blotter/api/pkg/tag"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	fp "github.com/OhYee/goutils/functional"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

var (
	htmlEscape, _     = regexp.Compile("<[^>]+>|\\s")
	charactorMatch, _ = regexp.Compile("[\u007f-\uffff]")
)

func CalcPostLength(html string) int {
	text := htmlEscape.ReplaceAllString(html, "")
	result := charactorMatch.FindAllString(text, -1)
	return len(result)
}

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
func QueryByURL(url string, post interface{}, statusPipeline bson.M, project bson.M) (total int64, err error) {
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
	if statusPipeline != nil {
		pipeline = append(pipeline, statusPipeline)
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
	cnt, err := QueryByURL(url, &postsDB, nil, nil)
	if cnt > 0 {
		post = postsDB[0].ToPost()
	}
	return
}

// GetPublicFieldPost get all field post
func GetPublicFieldPost(url string) (post PublicField, err error) {
	postsDB := make([]PublicFieldDB, 0)
	cnt, err := QueryByURL(url, &postsDB, bson.M{
		"$match": bson.M{"$or": []bson.M{
			{"status": 1}, // 隐藏
			{"status": 2}, // 发布
		}},
	}, nil)
	if cnt > 0 {
		post = postsDB[0].ToPost()
	}
	return
}

// IncView view +1
func IncView(url string) {
	if _, err := mongo.Update(
		"blotter", "posts", bson.M{"url": url},
		bson.M{"$inc": bson.M{"view": 1}}, nil,
	); err != nil {
		output.Err(err)
	}
}

func getPosts(
	status int8,
	offset int64, number int64,
	withTags []string,
	withoutTags []string,
	sortField string, sortType int,
	searchWord string, searchFields []string,
	res interface{},
) (total int64, err error) {
	if number < 0 {
		number = 0
	}

	sortQuery := bson.M{"$sort": bson.M{"publish_time": -1}}
	if sortField != "" && (sortType == 1 || sortType == -1) {
		sortQuery = bson.M{"$sort": bson.M{sortField: sortType}}
	}
	pipeline := []bson.M{}

	if status != -1 {
		pipeline = append(
			pipeline,
			bson.M{"$match": bson.M{"status": status}},
		)
	}

	withTagsID := mongo.StringToObjectIDs(withTags...)
	withoutTagsID := mongo.StringToObjectIDs(withoutTags...)
	if len(withTagsID) != 0 {
		pipeline = append(
			pipeline,
			bson.M{"$match": bson.M{"tags": bson.M{"$in": withTagsID}}},
		)
	}
	if len(withoutTagsID) != 0 {
		pipeline = append(
			pipeline,
			bson.M{"$match": bson.M{"tags": bson.M{"$nin": withoutTagsID}}},
		)
	}

	words := fp.FilterString(func(s string, idx int) bool {
		return len(strings.Replace(s, " ", "", -1)) != 0
	}, getJieba().CutForSearch(searchWord, true))

	if len(words) > 0 {
		s := make([]bson.M, len(searchFields)*len(words))
		wordsNumber := len(words)
		for i, ss := range searchFields {
			for j, word := range words {
				s[i*wordsNumber+j] = bson.M{ss: bson.M{"$regex": word, "$options": "i"}}
			}
		}

		pipeline = append(
			pipeline,
			bson.M{"$match": bson.M{"$or": s}},
		)
	}

	pipeline = append(
		pipeline,
		sortQuery,
	)

	pipeline = append(
		pipeline,
		bson.M{
			"$lookup": bson.M{
				"from":         "tags",
				"localField":   "tags",
				"foreignField": "_id",
				"as":           "tags",
			},
		},
	)

	if (offset != 0 || number != 0) && len(words) == 0 {
		pipeline = append(
			pipeline,
			bson.M{"$limit": offset + number},
			bson.M{"$skip": offset},
		)
	}

	posts := make([]SortPost, 0)
	total, err = Query(pipeline, &posts)

	if len(words) > 0 {
		Sort(posts, words, searchFields)
	}

	if (offset != 0 || number != 0) && len(words) != 0 {
		posts = posts[min(max(0, total-1), offset):min(total, offset+number)]
	}

	switch v := res.(type) {
	case *[]AdminFieldDB:
		(*v) = make([]AdminFieldDB, len(posts))
		for idx, p := range posts {
			(*v)[idx] = p.ToAdminDB()
		}
	case *[]CardField:
		(*v) = make([]CardField, len(posts))
		for idx, p := range posts {
			(*v)[idx] = p.ToCard()
		}
	}

	return
}

// GetCardPosts get all field post
func GetCardPosts(
	offset int64, number int64,
	withTags []string, withoutTags []string,
	sortField string, sortType int,
	searchWord string, searchFields []string,
	hidden bool,
) (total int64, posts []CardField, err error) {
	postsDB := make([]CardField, 0)
	if hidden {
		var hidden_tags []tag.Type
		hidden_tags, err = tag.GetHidden()
		if err != nil {
			return
		}
		for _, tag := range hidden_tags {
			withoutTags = append(withoutTags, tag.ID)
		}
	}
	total, err = getPosts(2, offset, number, withTags, withoutTags, sortField, sortType, searchWord, searchFields, &postsDB)

	posts = make([]CardField, len(postsDB))
	for idx, post := range postsDB {
		posts[idx] = post
	}
	return
}

// GetAdminPosts get all field post
func GetAdminPosts(
	offset int64, number int64,
	withTags []string, withoutTags []string,
	sortField string,
	sortType int,
	searchWord string, searchFields []string,
) (total int64, posts []AdminField, err error) {
	postsDB := make([]AdminFieldDB, 0)
	total, err = getPosts(-1, offset, number, withTags, withoutTags, sortField, sortType, searchWord, searchFields, &postsDB)

	posts = make([]AdminField, len(postsDB))
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
	status int8,
	// published bool,
	headImage string,
	images []string,
	poptext string,
) (err error) {

	html, err := markdown.Render(raw, true)
	if err != nil {
		return
	}
	tagsID := make([]primitive.ObjectID, len(tags))
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
	p.Status = status
	// p.Published = published
	p.HeadImage = headImage
	p.Images = images
	p.Length = int64(CalcPostLength(html))
	p.PopText = poptext

	if Existed(url) {
		err = fmt.Errorf("Post with url existed: %s", url)
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
	// published bool,
	status int8,
	headImage string,
	images []string,
	poptext string,
) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	html, err := markdown.Render(raw, true)
	if err != nil {
		return
	}
	tagsID := make([]primitive.ObjectID, len(tags))
	for idx, tagName := range tags {
		if tagsID[idx], err = primitive.ObjectIDFromHex(tagName); err != nil {
			return
		}
	}

	p := DB{}
	p.ID = objectID
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
	p.Status = status
	// p.Published = published
	p.HeadImage = headImage
	p.Images = images
	p.Length = int64(CalcPostLength(html))
	p.PopText = poptext

	_, err = mongo.Update("blotter", "posts", bson.M{
		"_id": objectID,
	}, bson.M{
		"$set": p,
	}, nil)
	return
}

func Delete(id string) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err == nil {
		mongo.Remove("blotter", "posts", bson.M{
			"_id": objectID,
		}, nil)
	}
}
