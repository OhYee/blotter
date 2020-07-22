package friends

// Simple friend type
type Simple struct {
	Link string `json:"link" bson:"link"`
	Name string `json:"name" bson:"name"`
}

// Friend friend type
type Friend struct {
	Simple      `bson:",inline"`
	Image       string       `json:"image" bson:"image"`
	Description string       `json:"description" bson:"description"`
	RSS         string       `json:"rss" bson:"rss"`
	Posts       []FriendPost `json:"posts" bson:"posts"`
	Error       bool         `json:"error" bson:"error"`
}

// FriendPost post of friend
type FriendPost struct {
	Title string `json:"title" bson:"title"`
	Link  string `json:"link" bson:"link"`
}

// WithIndex friend type with index
type WithIndex struct {
	Friend `bson:",inline"`
	Index  int `json:"index" bson:"index"`
}
