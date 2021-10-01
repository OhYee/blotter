package comment

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DatabaseName   = "blotter"
	CollectionName = "comments"
)

var (
	ErrShake        = fmt.Errorf("anti-shake")
	defaultObjectID = primitive.ObjectID{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
)
