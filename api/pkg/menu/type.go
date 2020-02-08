package menu

// Type of the blotter
type Type struct {
	Icon string `json:"icon" bson:"icon"`
	Name string `json:"name" bson:"name"`
	Link string `json:"link" bson:"link"`
}
