package tag

type Base struct {
	ID    string `json:"id" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Short string `json:"short" bson:"short"`
	Icon  string `json:"icon" bson:"icon"`
	Color string `json:"color" bson:"color"`
}

// Type type
type Type struct {
	Base        `bson:",inline"`
	Description string `json:"description" bson:"description"`
}

// WithCount tag type with count
type WithCount struct {
	Base  `bson:",inline"`
	Count int64 `json:"count" bson:"count"`
}
