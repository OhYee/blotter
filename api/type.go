package api

// Menu of the blotter
type Menu struct {
	Icon string `json:"icon"`
	Name string `json:"name"`
	Link string `json:"link"`
}

type Friend struct {
	Image       string            `json:"image"`
	Link        string            `json:"link"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Posts       []FriendPost `json:"posts"`
}


type FriendPost struct {
	Title string `json:"title"`
	Link  string `json:"link"`
}

