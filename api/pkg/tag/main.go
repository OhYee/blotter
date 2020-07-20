package tag

import (
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetTags get all tags with count
func GetTags(keyword string, offset int64, number int64, sortField string, sortInc bool) (total int64, res []WithCount, err error) {
	res = make([]WithCount, 0)

	pipeline := []bson.M{}
	if keyword != "" {
		pipeline = append(
			pipeline,
			bson.M{
				"$match": bson.M{
					"$or": []bson.M{
						{"name": bson.M{"$regex": keyword}},
						{"short": bson.M{"$regex": keyword}},
					},
				},
			},
		)
	}

	if sortField == "" {
		sortField = "count"
		sortInc = false
	}

	output.Debug("%v %v", sortField, sortInc)
	pipeline = append(
		pipeline,
		bson.M{
			"$lookup": bson.M{
				"from":         "posts",
				"localField":   "_id",
				"foreignField": "tags",
				"as":           "posts",
			},
		},
		bson.M{
			"$set": bson.M{"count": bson.M{"$size": "$posts"}},
		},
		bson.M{
			"$sort": bson.M{sortField: map[bool]int{true: 1, false: -1}[sortInc]},
		},
	)

	if number != 0 {
		pipeline = append(pipeline, mongo.AggregateOffset(offset, number)...)
	}

	total, err = mongo.Aggregate("blotter", "tags", pipeline, nil, &res)
	return
}

// New tag
func New(name string, short string, color string, icon string, description string) (err error) {
	exist, err := Existed(primitive.NilObjectID.Hex(), short)
	if err != nil {
		return
	}
	if exist {
		err = errors.New("Tat %s has existed", short)
		return
	}

	// New tag
	_, err = mongo.Add(
		"blotter",
		"tags",
		nil,
		bson.M{
			"name":        name,
			"short":       short,
			"color":       color,
			"icon":        icon,
			"description": description,
		},
	)
	return
}

// Update tag data
func Update(id string, name string, short string, color string, icon string, description string) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	existed, err := Existed(id, short)
	if err != nil {
		return
	}
	if existed {
		err = errors.New("Tat %s has existed", short)
		return
	}

	// Update all posts with tag
	_, err = mongo.Update(
		"blotter", "tags",
		bson.M{"_id": objectID},
		bson.M{
			"$set": bson.M{
				"name":        name,
				"short":       short,
				"color":       color,
				"icon":        icon,
				"description": description,
			},
		},
		nil,
	)
	return
}

// Delete tag
func Delete(id string) (err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	// Delete tag
	_, err = mongo.Remove("blotter", "tags", bson.M{
		"_id": objectID,
	}, nil)

	if err != nil {
		return
	}

	// Delete all posts with tag
	_, err = mongo.Update(
		"blotter", "posts",
		bson.M{"tags": objectID},
		bson.M{
			"$pull": bson.M{"tags": objectID},
		},
		nil,
	)
	return
}

// Existed tag(short) has existed
func Existed(id string, short string) (existed bool, err error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	count, err := mongo.Find("blotter", "tags", bson.M{
		"short": short,
		"_id":   bson.M{"$ne": objectID},
	}, nil, nil)

	existed = count != 0
	return
}

// Get tag with short
func Get(short string) (tag Type, err error) {
	tags := make([]Type, 0)

	total, err := mongo.Find("blotter", "tags", bson.M{
		"short": short,
	}, nil, &tags)

	if err != nil {
		return
	}
	if total < 1 {
		err = errors.New("No tag %s", short)
		return
	}
	tag = tags[0]
	return
}
