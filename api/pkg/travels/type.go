package travels

// Type travels city
type Type struct {
	Name    string  `json:"name" bson:"name"`
	Lng     float64 `json:"lng" bson:"lng"`
	Lat     float64 `json:"lat" bson:"lat"`
	Zoom    float64 `json:"zoom" bson:"zoom"`
	Travels []struct {
		Time uint64 `json:"time" bson:"time"`
		Link string `json:"link" bson:"link"`
	} `json:"travels" bson:"travels"`
}

// Type travels city
type Travel struct {
	Name string  `json:"name" bson:"name"`
	Lng  float64 `json:"lng" bson:"lng"`
	Lat  float64 `json:"lat" bson:"lat"`
	Zoom float64 `json:"zoom" bson:"zoom"`
	Time uint64  `json:"time" bson:"time"`
	Link string  `json:"link" bson:"link"`
}
