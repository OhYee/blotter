package api

// Menu of the blotter
type Menu struct {
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

type PostCard struct {
	Title       string `json:"title" bson:"title"`
	Abstract    string `json:"abstract" bson:"abstract"`
	View        int    `json:"view" bson:"view"`
	URL         string `json:"url" bson:"url"`
	PublishTime int64  `json:"publish_time" bson:"publish_time"`
	EditTime    int64  `json:"edit_time" bson:"edit_time"`
	Tags        []Tag  `json:"tags" bson:"tags"`
	HeadImage   string `json:"head_image" bson:"head_image"`
}

type Tag struct {
	Name  string `json:"name" bson:"name"`
	Short string `json:"short" bson:"short"`
	Icon  string `json:"icon" bson:"icon"`
	Color string `json:"color" bson:"color"`
}
