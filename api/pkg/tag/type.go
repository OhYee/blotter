package tag

// Type type
type Type struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Short string `json:"short" bson:"short"`
	Icon  string `json:"icon" bson:"icon"`
	Color string `json:"color" bson:"color"`
}

// WithCount tag type with count
type WithCount struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Short string `json:"short" bson:"short"`
	Icon  string `json:"icon" bson:"icon"`
	Color string `json:"color" bson:"color"`
	Count int64  `json:"count" bson:"count"`
}
