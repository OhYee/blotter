package friends

// Friend friend type
type Friend struct {
	Image       string       `json:"image" bson:"image"`
	Link        string       `json:"link" bson:"link"`
	Name        string       `json:"name" bson:"name"`
	Description string       `json:"description" bson:"description"`
	Posts       []FriendPost `json:"posts" bson:"posts"`
}

// FriendPost post of friend
type FriendPost struct {
	Title string `json:"title" bson:"title"`
	Link  string `json:"link" bson:"link"`
}
