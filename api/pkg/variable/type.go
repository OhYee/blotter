package variable

import (
	"encoding/json"

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

// SetBool set value of key to bool
func (v Variables) SetBool(key string, value *bool, defaultValue bool) (err error) {
	var t interface{}
	var ok bool
	if t, ok = v[key]; !ok {
		output.Err(errors.New("Can not get value of %s", key))
		*value = defaultValue
		return
	}
	if *value, ok = t.(bool); !ok {
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
	switch value.(type) {
	case string:
		s = value.(string)
	default:
		s = ""
	}
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

// Experience for about api
type Experience struct {
	Name  string `json:"name"`
	Major string `json:"major"`
	Time  string `json:"time"`
}

// Award for about api
type Award struct {
	Name  string `json:"name"`
	Level string `json:"level"`
	Count int64  `json:"count"`
}
type BlotterVariables struct {
	ADInner           string       `json:"ad_inner" bson:"ad_inner"`
	ADShow            string       `json:"ad_show" bson:"ad_show"`
	ADText            string       `json:"ad_text" bson:"ad_text"`
	Author            string       `json:"author" bson:"author"`
	Avatar            string       `json:"avatar" bson:"avatar"`
	Awards            []Award      `json:"awards" bson:"awards"`
	Beian             string       `json:"beian" bson:"beian"`
	BlogName          string       `json:"blog_name" bson:"blog_name"`
	Description       string       `json:"description" bson:"description"`
	Edu               []Experience `json:"edu" bson:"edu"`
	Email             string       `json:"email" bson:"email"`
	From              string       `json:"from" bson:"from"`
	Github            string       `json:"github" bson:"github"`
	GithubID          string       `json:"github_id" bson:"github_id"`
	GithubRedirect    string       `json:"github_redirect" bson:"github_redirect"`
	GithubSecret      string       `json:"github_secret" bson:"github_secret"`
	Grey              bool         `json:"grey" bson:"grey"`
	Head              string       `json:"head" bson:"head"`
	QiniuAccessKey    string       `json:"qiniu_access_key" bson:"qiniu_access_key"`
	QiniuPrefix       string       `json:"qiniu_prefix" bson:"qiniu_prefix"`
	QiniuSecretKey    string       `json:"qiniu_secret_key" bson:"qiniu_secret_key"`
	QiniuStaticDomain string       `json:"qiniu_static_domain" bson:"qiniu_static_domain"`
	QQ                string       `json:"qq" bson:"qq"`
	QQID              string       `json:"qq_id" bson:"qq_id"`
	QQKey             string       `json:"qq_key" bson:"qq_key"`
	QQRedirect        string       `json:"qq_redirect" bson:"qq_redirect"`
	QQRobot           string       `json:"qqrobot" bson:"qqrobot"`
	Quote             string       `json:"quote" bson:"quote"`
	Root              string       `json:"root" bson:"root"`
	SMTPAddress       string       `json:"smtp_address" bson:"smtp_address"`
	SMTPPassword      string       `json:"smtp_password" bson:"smtp_password"`
	SMTPUser          string       `json:"smtp_user" bson:"smtp_user"`
	SMTPUsername      string       `json:"smtp_username" bson:"smtp_username"`
	View              int          `json:"view" bson:"view"`
	Vmess             string       `json:"vmess" bson:"vmess"`
	Zhihu             string       `json:"zhihu" bson:"zhihu"`
}

func NewBlotterVariables(vars []map[string]interface{}) (res BlotterVariables, err error) {
	m := make(map[string]interface{})
	for _, v := range vars {
		m[v["key"].(string)] = v["value"]
	}
	output.Debug("%+v", m)
	b, err := json.Marshal(m)
	if err != nil {
		return
	}

	err = json.Unmarshal(b, &res)
	return
}
