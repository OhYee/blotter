package mongo

import (
	"context"

	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions = options.Client().ApplyURI("mongodb://127.0.0.1:27017")

type countResult struct {
	Count int64 `bson:"count"`
}

type Conn struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

func NewConn(databaseName string, collectionName string) (conn *Conn, err error) {
	defer func() {
		if err != nil {
			err = errors.NewErr(err)
		}
	}()
	conn = &Conn{}

	conn.Client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		return
	}

	err = conn.Client.Ping(context.TODO(), nil)
	if err != nil {
		return
	}

	conn.Collection = conn.Client.Database(databaseName).Collection(collectionName)
	return
}

func (conn *Conn) Close() {
	if conn.Client != nil {
		conn.Client.Disconnect(context.TODO())
	}
}

func Find(databaseName string, collectionName string, filter interface{},
	opt *options.FindOptions, res interface{}) (total int64, err error) {
	defer func() {
		if err != nil {
			err = errors.NewErr(err)
		}
	}()

	conn, err := NewConn(databaseName, collectionName)
	defer conn.Close()
	if err != nil {
		return
	}

	cur, err := conn.Collection.Find(context.TODO(), filter, opt)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	if total, err = conn.Collection.CountDocuments(context.TODO(), filter, nil); err != nil {
		return
	}

	if res != nil {
		err = cur.All(context.TODO(), res)
	}
	return
}

func Aggregate(databaseName string, collectionName string, pipeline interface{},
	opt *options.AggregateOptions, res interface{}) (total int64, err error) {
	defer func() {
		if err != nil {
			err = errors.NewErr(err)
		}
	}()

	conn, err := NewConn(databaseName, collectionName)
	defer conn.Close()
	if err != nil {
		return
	}

	cur, err := conn.Collection.Aggregate(context.TODO(), pipeline, opt)
	if err != nil {
		return
	}
	defer cur.Close(context.TODO())

	if res != nil {
		if err = cur.All(context.TODO(), res); err != nil {
			return
		}
	}

	count := countResult{}
	countPipeline, err := pipelineTruncated(pipeline)
	countPipeline = append(countPipeline, bson.M{"$count": "count"})
	if err != nil {
		return
	}
	countCur, err := conn.Collection.Aggregate(context.TODO(), countPipeline, opt)
	if err != nil {
		return
	}
	defer countCur.Close(context.TODO())
	if countCur.Next(context.TODO()) {
		if err = countCur.Decode(&count); err != nil {
			return
		}
	}
	total = count.Count

	return
}

func bsonFormat(b interface{}) (bb []bson.M, err error) {
	switch b.(type) {
	case bson.D:
		bb = []bson.M{b.(bson.D).Map()}
	case []bson.E:
		bb = []bson.M{bson.D(b.([]bson.E)).Map()}
	case bson.E:
		bb = []bson.M{bson.D([]bson.E{b.(bson.E)}).Map()}
	case bson.M:
		bb = []bson.M{b.(bson.M)}
	case map[string]interface{}:
		bb = []bson.M{bson.M(b.(map[string]interface{}))}
	case []bson.M:
		bb = b.([]bson.M)
	case []map[string]interface{}:
		m := b.([]map[string]interface{})
		bb = make([]bson.M, len(m))
		for idx, data := range m {
			bb[idx] = bson.M(data)
		}
	default:
		err = errors.New("Can format bson: %+v", b)
		bb = []bson.M{}
	}
	return
}

func pipelineTruncated(pipeline interface{}) (res []bson.M, err error) {
	m, err := bsonFormat(pipeline)
	end := -1
	for i := len(m) - 1; i >= 0; i-- {
		if _, exist := m[i]["$limit"]; exist {
			continue
		}
		if _, exist := m[i]["$skip"]; exist {
			continue
		}
		end = i
		break
	}
	res = m[0 : end+1]
	return
}

func Add(databaseName string, collectionName string,
	opt *options.InsertManyOptions, documents ...interface{}) (ids []interface{}, err error) {
	conn, err := NewConn(databaseName, collectionName)
	defer conn.Close()
	if err != nil {
		return
	}

	result, err := conn.Collection.InsertMany(context.TODO(), documents, opt)
	if err != nil {
		return
	}
	ids = result.InsertedIDs
	return
}

func Update(databaseName string, collectionName string, filter interface{}, update interface{},
	opt *options.UpdateOptions) (result *mongo.UpdateResult, err error) {
	conn, err := NewConn(databaseName, collectionName)
	defer conn.Close()
	if err != nil {
		return
	}

	result, err = conn.Collection.UpdateMany(context.TODO(), filter, update, opt)
	if err != nil {
		return
	}
	return
}

func Remove(databaseName string, collectionName string, filter interface{},
	opt *options.DeleteOptions) (count int64, err error) {
	conn, err := NewConn(databaseName, collectionName)
	defer conn.Close()
	if err != nil {
		return
	}

	result, err := conn.Collection.DeleteMany(context.TODO(), filter, opt)
	if err != nil {
		return
	}
	count = result.DeletedCount
	return
}

// AggregateOffset using offset in aggregate
func AggregateOffset(offset int64, number int64) []bson.M {
	return []bson.M{
		bson.M{"$limit": offset + number},
		bson.M{"$skip": offset},
	}
}
