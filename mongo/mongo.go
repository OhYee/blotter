package mongo

import (
	"context"
	"github.com/OhYee/rainbow/errors"
	"github.com/OhYee/blotter/output"

	// "go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions = options.Client().ApplyURI("mongodb://127.0.0.1:27017")

func Query(databaseName string, collectionName string, filter interface{}, res interface{}) (err error) {
	defer func() {
		if err != nil {
			err = errors.NewErr(err)
		}
	}()

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}
	defer client.Disconnect(context.TODO())

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		return
	}

	collection := client.Database(databaseName).Collection(collectionName)
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	cur.All(context.TODO(), res)
	output.Debug("query result: %+v",res)
	return
}
