package post

import (
	"github.com/OhYee/blotter/api/pkg/tag"
	"github.com/OhYee/goutils/time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TimeDB post time in db (int64)
type TimeDB struct {
	PublishTime int64 `json:"publish_time" bson:"publish_time"`
	EditTime    int64 `json:"edit_time" bson:"edit_time"`
}

// ToTime transfer TimeDB to Time
func (t TimeDB) ToTime() Time {
	return Time{
		PublishTime: time.ToString(t.PublishTime),
		EditTime:    time.ToString(t.EditTime),
	}
}

// Time post time (string)
type Time struct {
	PublishTime string `json:"publish_time" bson:"publish_time"`
	EditTime    string `json:"edit_time" bson:"edit_time"`
}

// ToTimeDB transfer Time to TimeDB
func (t Time) ToTimeDB() TimeDB {
	return TimeDB{
		PublishTime: time.FromString(t.PublishTime),
		EditTime:    time.FromString(t.EditTime),
	}
}

// CardProps props of post card
type CardProps struct {
	Time      `bson:",inline"`
	Title     string `json:"title" bson:"title"`
	Abstract  string `json:"abstract" bson:"abstract"`
	View      int    `json:"view" bson:"view"`
	URL       string `json:"url" bson:"url"`
	HeadImage string `json:"head_image" bson:"head_image"`
}

// ToCardDBProps transfer CardProps to CardDBProps
func (props CardProps) ToCardDBProps() CardDBProps {
	return CardDBProps{
		Title:     props.Title,
		Abstract:  props.Abstract,
		View:      props.View,
		URL:       props.URL,
		HeadImage: props.HeadImage,
		TimeDB:    props.Time.ToTimeDB(),
	}
}

// CardDBProps props of post card database type
type CardDBProps struct {
	TimeDB    `bson:",inline"`
	Title     string `json:"title" bson:"title"`
	Abstract  string `json:"abstract" bson:"abstract"`
	View      int    `json:"view" bson:"view"`
	URL       string `json:"url" bson:"url"`
	HeadImage string `json:"head_image" bson:"head_image"`
}

// ToCardProps transfer CardDBProps to CardProps
func (props CardDBProps) ToCardProps() CardProps {
	return CardProps{
		Title:     props.Title,
		Abstract:  props.Abstract,
		View:      props.View,
		URL:       props.URL,
		HeadImage: props.HeadImage,
		Time:      props.TimeDB.ToTime(),
	}
}

// PublicProps extra props of public post
type PublicProps struct {
	Content string `json:"content" bson:"content"`
}

// EditProps extra props for editing
type EditProps struct {
	Raw       string   `json:"raw" bson:"raw"`
	Keywords  []string `json:"keywords" bson:"keywords"`
	Published bool     `json:"published" bson:"published"`
}

/*

	PostCard

*/

// CardFieldDB PostCard type in database
type CardFieldDB struct {
	CardDBProps `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToCard transfer PostCardDB to PostCard
func (post CardFieldDB) ToCard() CardField {
	return CardField{
		CardProps: post.CardDBProps.ToCardProps(),
		Tags:      post.Tags,
	}
}

// CardField PostCard type
type CardField struct {
	CardProps `bson:",inline"`
	Tags      []tag.Type `json:"tags" bson:"tags"`
}

// ToCardDB transfer PostCard to PostCardDB
func (post CardField) ToCardDB() CardFieldDB {
	return CardFieldDB{
		CardDBProps: post.CardProps.ToCardDBProps(),
		Tags:        post.Tags,
	}
}

/*

	Post

*/

// PublicFieldDB post type in database
type PublicFieldDB struct {
	CardDBProps `bson:",inline"`
	PublicProps `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToPost transfer PostDB to Post
func (post PublicFieldDB) ToPost() PublicField {
	return PublicField{
		CardProps:   post.CardDBProps.ToCardProps(),
		PublicProps: post.PublicProps,
		Tags:        post.Tags,
	}
}

// PublicField post type for show
type PublicField struct {
	CardProps   `bson:",inline"`
	PublicProps `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToDB transfer Post to ToPostDB
func (post PublicField) ToDB() PublicFieldDB {
	return PublicFieldDB{
		CardDBProps: post.CardProps.ToCardDBProps(),
		PublicProps: post.PublicProps,
		Tags:        post.Tags,
	}
}

/*

	PostAll

*/

// CompleteFieldDB post with all field for editing in database
type CompleteFieldDB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CardDBProps `bson:",inline"`
	PublicProps `bson:",inline"`
	EditProps   `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToPost transfer CompleteFieldDB to CompleteField
func (post CompleteFieldDB) ToPost() CompleteField {
	return CompleteField{
		ID:          post.ID.Hex(),
		CardProps:   post.CardDBProps.ToCardProps(),
		PublicProps: post.PublicProps,
		EditProps:   post.EditProps,
		Tags:        post.Tags,
	}
}

// CompleteField post with all field for editing
type CompleteField struct {
	ID          string `json:"id" bson:"_id"`
	CardProps   `bson:",inline"`
	PublicProps `bson:",inline"`
	EditProps   `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}
