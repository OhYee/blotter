package api

// MenuItem of the blotter
type MenuItem struct {
	Icon string `json:"icon" bson:"icon"`
	Name string `json:"name" bson:"name"`
	Link string `json:"link" bson:"link"`
}

type Friend struct {
	Image       string       `json:"image" bson:"image"`
	Link        string       `json:"link" bson:"link"`
	Name        string       `json:"name" bson:"name"`
	Description string       `json:"description" bson:"description"`
	Posts       []FriendPost `json:"posts" bson:"posts"`
}

type FriendPost struct {
	Title string `json:"title" bson:"title"`
	Link  string `json:"link" bson:"link"`
}

type PostUnix struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime int64  `json:"publish_time" bson:"publish_time"`
	EditTime    int64  `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
	Content     string `json:"content" bson:"content"`
}

type PostTime struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime string `json:"publish_time" bson:"publish_time"`
	EditTime    string `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
	Content     string `json:"content" bson:"content"`
}

type PostDate struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime int64  `json:"publish_time" bson:"publish_time"`
	EditTime    int64  `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
	Content     string `json:"content" bson:"content"`
}

type PostCardUnix struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime int64  `json:"publish_time" bson:"publish_time"`
	EditTime    int64  `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
}

type PostCardTime struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime string `json:"publish_time" bson:"publish_time"`
	EditTime    string `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
}

type Tag struct {
	Name  string `json:"name" bson:"name"`
	Short string `json:"short" bson:"short"`
	Icon  string `json:"icon" bson:"icon"`
	Color string `json:"color" bson:"color"`
}

type CommentUnix struct {
	ID      int    `json:"id" bson:"id"`
	Mail    string `json:"email" bson:"email"`
	Avatar  string `json:"avatar" bson:"avatar"`
	Time    int64  `json:"time" bson:"time"`
	Content string `json:"content" bson:"content"`
	Reply   int    `json:"reply" bson:"reply"`
}

type CommentTime struct {
	ID       int            `json:"id" bson:"id"`
	Mail     string         `json:"email" bson:"email"`
	Avatar   string         `json:"avatar" bson:"avatar"`
	Time     string         `json:"time" bson:"time"`
	Content  string         `json:"content" bson:"content"`
	Children []*CommentTime `json:"children" bson:"children"`
}

type Variable struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
}
