package variable

import (
	"github.com/OhYee/rainbow/errors"
)

// Variables type
type Variables map[string]interface{}

// SetString set value of key to string
func (v Variables) SetString(key string, value *string) (err error) {
	var t interface{}
	var ok bool
	if t, ok = v[key]; !ok {
		err = errors.New("Can not get value of %s", key)
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
