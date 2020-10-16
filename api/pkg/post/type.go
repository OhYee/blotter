package post

import (
	"sort"
	"strings"

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
	PublishTime int64  `json:"publish_time" bson:"publish_time"`
	EditTime    int64  `json:"edit_time" bson:"edit_time"`
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int64  `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	HeadImage   string `json:"head_image" bson:"head_image"`
	Length      int64  `json:"length" bson:"length"`
}

// PublicProps extra props of public post
type PublicProps struct {
	Content string   `json:"content" bson:"content"`
	Images  []string `json:"images" bson:"images"`
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

// AdminField PostCard type
type AdminField struct {
	CardProps `bson:",inline"`
	ID        string     `json:"id" bson:"_id"`
	Published bool       `json:"published" bson:"published"`
	Tags      []tag.Type `json:"tags" bson:"tags"`
}

// ToCardDB transfer PostCard to PostCardDB
func (post AdminField) ToCardDB() AdminFieldDB {
	id, err := primitive.ObjectIDFromHex(post.ID)
	if err != nil {
		id = primitive.ObjectID{}
	}
	return AdminFieldDB{
		ID:        id,
		CardProps: post.CardProps,
		Published: post.Published,
		Tags:      post.Tags,
	}
}

// AdminFieldDB PostCard type in database
type AdminFieldDB struct {
	CardProps `bson:",inline"`
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	Published bool               `json:"published" bson:"published"`
	Tags      []tag.Type         `json:"tags" bson:"tags"`
}

// ToCard transfer PostCardDB to PostCard
func (post AdminFieldDB) ToCard() AdminField {
	return AdminField{
		ID:        post.ID.Hex(),
		CardProps: post.CardProps,
		Published: post.Published,
		Tags:      post.Tags,
	}
}

// CardField PostCard type
type CardField struct {
	CardProps `bson:",inline"`
	Tags      []tag.Type `json:"tags" bson:"tags"`
}

/*

	Post

*/

// PublicFieldDB post type in database
type PublicFieldDB struct {
	CardProps   `bson:",inline"`
	PublicProps `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToPost transfer PostDB to Post
func (post PublicFieldDB) ToPost() PublicField {
	return PublicField{
		CardProps:   post.CardProps,
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
		CardProps:   post.CardProps,
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
	CardProps   `bson:",inline"`
	PublicProps `bson:",inline"`
	EditProps   `bson:",inline"`
	Tags        []tag.Type `json:"tags" bson:"tags"`
}

// ToPost transfer CompleteFieldDB to CompleteField
func (post CompleteFieldDB) ToPost() CompleteField {
	return CompleteField{
		ID:          post.ID.Hex(),
		CardProps:   post.CardProps,
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

// DB type in database posts collection
type DB struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CardProps   `bson:",inline"`
	PublicProps `bson:",inline"`
	EditProps   `bson:",inline"`
	Tags        []primitive.ObjectID `json:"tags" bson:"tags"`
}

type SortPost CompleteFieldDB

func (post SortPost) ToCard() CardField {
	return CardField{
		CardProps: post.CardProps,
		Tags:      post.Tags,
	}
}

func (post SortPost) ToAdminDB() AdminFieldDB {
	return AdminFieldDB{
		CardProps: post.CardProps,
		ID:        post.ID,
		Published: post.Published,
		Tags:      post.Tags,
	}
}

type SortPostArray struct {
	Array  []SortPost
	scores []int
}

func Sort(posts []SortPost, words []string, fields []string) {
	scores := make([]int, len(posts))
	for idx, post := range posts {
		scores[idx] = 0
		for _, field := range fields {
			value := ""
			switch field {
			case "title":
				value = strings.ToLower(post.Title)
			case "abstract":
				value = strings.ToLower(post.Abstract)
			case "raw":
				value = strings.ToLower(post.Raw)
			}
			for _, word := range words {
				if strings.Index(value, strings.ToLower(word)) >= 0 {
					switch field {
					case "title":
						scores[idx] += 4
					case "abstract":
						scores[idx] += 2
					case "raw":
						scores[idx]++
					}
				}
			}
		}
	}

	sort.Sort(SortPostArray{
		Array:  posts,
		scores: scores,
	})
}

func (arr SortPostArray) Len() int {
	return len(arr.Array)
}

func (arr SortPostArray) Less(i, j int) bool {
	return arr.scores[i] > arr.scores[j]
}

func (arr SortPostArray) Swap(i, j int) {
	arr.Array[i], arr.Array[j] = arr.Array[j], arr.Array[i]
	arr.scores[i], arr.scores[j] = arr.scores[j], arr.scores[i]
}
