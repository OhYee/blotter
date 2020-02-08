package post

import (
	"github.com/OhYee/blotter/api/pkg/tag"
	"github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*

	PostAll

*/

// CompleteFieldDB post with all field for editing in database
type CompleteFieldDB struct {
	ID          primitive.ObjectID   `json:"id" bson:"_id"`
	Title       string               `json:"title" bson:"title"`
	Abstract    string               `json:"abstract" bson:"abstract"`
	View        int64                `json:"view" bson:"view"`
	URL         string               `json:"url" bson:"url"`
	PublishTime int64                `json:"publish_time" bson:"publish_time"`
	EditTime    int64                `json:"edit_time" bson:"edit_time"`
	Content     string               `json:"content" bson:"content"`
	Raw         string               `json:"raw" bson:"raw"`
	Tags        []primitive.ObjectID `json:"tags" bson:"tags"`
	Keywords    []string             `json:"keywords" bson:"keywords"`
	Published   bool                 `json:"published" bson:"published"`
	HeadImage   string               `json:"head_image" bson:"head_image"`
}

// ToPost transfer CompleteFieldDB to CompleteField
func (post CompleteFieldDB) ToPost() CompleteField {
	return CompleteField{
		ID:          post.ID.Hex(),
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: post.PublishTime,
		EditTime:    post.EditTime,
		Content:     post.Content,
		Raw:         post.Raw,
		// Tags:        post.Tags,
		Keywords:  post.Keywords,
		Published: post.Published,
		HeadImage: post.HeadImage,
	}
}

// CompleteField post with all field for editing
type CompleteField struct {
	ID          string     `json:"id" bson:"_id"`
	Title       string     `json:"title" bson:"title"`
	Abstract    string     `json:"abstract" bson:"abstract"`
	View        int64      `json:"view" bson:"view"`
	URL         string     `json:"url" bson:"url"`
	PublishTime int64      `json:"publish_time" bson:"publish_time"`
	EditTime    int64      `json:"edit_time" bson:"edit_time"`
	Content     string     `json:"content" bson:"content"`
	Raw         string     `json:"raw" bson:"raw"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
	Keywords    []string   `json:"keywords" bson:"keywords"`
	Published   bool       `json:"published" bson:"published"`
	HeadImage   string     `json:"head_image" bson:"head_image"`
}

/*

	Post

*/

// PublicFieldDB post type in database
type PublicFieldDB struct {
	Title       string     `json:"title" bson:"title"`
	Abstract    string     `json:"abstract" bson:"abstract"`
	View        int        `json:"view" bson:"view"`
	URL         string     `json:"url" bson:"url"`
	PublishTime int64      `json:"publish_time" bson:"publish_time"`
	EditTime    int64      `json:"edit_time" bson:"edit_time"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
	HeadImage   string     `json:"head_image" bson:"head_image"`
	Content     string     `json:"content" bson:"content"`
}

// ToPost transfer PostDB to Post
func (post PublicFieldDB) ToPost() PublicField {
	return PublicField{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.ToString(post.PublishTime),
		EditTime:    time.ToString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
		Content:     post.Content,
	}
}

// PublicField post type for show
type PublicField struct {
	Title       string     `json:"title" bson:"title"`
	Abstract    string     `json:"abstract" bson:"abstract"`
	View        int        `json:"view" bson:"view"`
	URL         string     `json:"url" bson:"url"`
	PublishTime string     `json:"publish_time" bson:"publish_time"`
	EditTime    string     `json:"edit_time" bson:"edit_time"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
	HeadImage   string     `json:"head_image" bson:"head_image"`
	Content     string     `json:"content" bson:"content"`
}

// ToDB transfer Post to ToPostDB
func (post PublicField) ToDB() PublicFieldDB {
	return PublicFieldDB{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.FromString(post.PublishTime),
		EditTime:    time.FromString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
		Content:     post.Content,
	}
}

/*

	PostCard

*/

// CardFieldDB PostCard type in database
type CardFieldDB struct {
	Title       string     `json:"title" bson:"title"`
	Abstract    string     `json:"abstract" bson:"abstract"`
	View        int        `json:"view" bson:"view"`
	URL         string     `json:"url" bson:"url"`
	PublishTime int64      `json:"publish_time" bson:"publish_time"`
	EditTime    int64      `json:"edit_time" bson:"edit_time"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
	HeadImage   string     `json:"head_image" bson:"head_image"`
}

// ToCard transfer PostCardDB to PostCard
func (post CardFieldDB) ToCard() CardField {
	return CardField{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.ToString(post.PublishTime),
		EditTime:    time.ToString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
	}
}

// CardField PostCard type
type CardField struct {
	Title       string     `json:"title" bson:"title"`
	Abstract    string     `json:"abstract" bson:"abstract"`
	View        int        `json:"view" bson:"view"`
	URL         string     `json:"url" bson:"url"`
	PublishTime string     `json:"publish_time" bson:"publish_time"`
	EditTime    string     `json:"edit_time" bson:"edit_time"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
	HeadImage   string     `json:"head_image" bson:"head_image"`
}

// ToCardDB transfer PostCard to PostCardDB
func (post CardField) ToCardDB() CardFieldDB {
	return CardFieldDB{
		Title:       post.Title,
		Abstract:    post.Abstract,
		View:        post.View,
		URL:         post.URL,
		PublishTime: time.FromString(post.PublishTime),
		EditTime:    time.FromString(post.EditTime),
		Tags:        post.Tags,
		HeadImage:   post.HeadImage,
	}
}
