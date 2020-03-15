package variable

import (
	"github.com/mitchellh/mapstructure"

	"github.com/OhYee/blotter/output"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Variables type
type Variables map[string]interface{}

// SetString set value of key to string
func (v Variables) SetString(key string, value *string) (err error) {
	var t interface{}
	var ok bool
	if t, ok = v[key]; !ok {
		output.Err(errors.New("Can not get value of %s", key))
		*value = ""
		return
	}
	if *value, ok = t.(string); !ok {
		err = errors.New("Value of %s is %s %T, not %t", key, t, *value)
		return
	}
	return err
}

// SetInt64 set value of key to int64
func (v Variables) SetInt64(key string, value *int64) (err error) {
	var t interface{}
	var ok bool
	if t, ok = v[key]; !ok {
		err = errors.New("Can not get value of %s", key)
		return
	}
	if *value, ok = t.(int64); !ok {
		err = errors.New("Value of %s is %s %T, not %t", key, t, *value)
		return
	}
	return err
}

// SetArray set value of key to int64
func (v Variables) SetArray(key string, value interface{}) (err error) {
	var t interface{}
	var ok bool
	var array primitive.A
	if t, ok = v[key]; !ok {
		err = errors.New("Can not get value of %s", key)
		return
	}
	if array, ok = t.(primitive.A); !ok {
		err = errors.New("Value of %s is %s %T, not %t", key, t, t, array)
		return
	}
	m := make([]map[string]interface{}, len(array))
	for idx, item := range array {
		switch item.(type) {
		case primitive.D:
			m[idx] = item.(primitive.D).Map()
		case map[string]interface{}:
			m[idx] = item.(map[string]interface{})
		case primitive.M:
			m[idx] = item.(primitive.M)
		default:
			err = errors.New("Can not transfer %s array item  %+v %T to map[string]interface{}", key, item, item)
			return
		}
	}
	err = mapstructure.Decode(m, value)
	return
}

// GetString get string value
func (v Variables) GetString(key string) (s string, exist bool) {
	value, exist := v[key]
	s = value.(string)
	return
}

// Type variable type
type Type struct {
	Key   string      `json:"key" bson:"key"`
	Value interface{} `json:"value" bson:"value"`
}

// FromMapToType transfer map[string]interface{} to Type
func FromMapToType(m map[string]interface{}) (t Type) {
	return Type{
		Key:   m["key"].(string),
		Value: m["value"],
	}
}

// FromMapSliceToTypeSlice transfer []map[string]interface{} to []Type
func FromMapSliceToTypeSlice(ms []map[string]interface{}) (ts []Type) {
	ts = make([]Type, len(ms))
	for idx, m := range ms {
		ts[idx] = FromMapToType(m)
	}
	return
}
