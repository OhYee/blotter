package variable

import (
	"fmt"
	"strings"

	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/goutils/transfer"
	"github.com/OhYee/rainbow/errors"
	"go.mongodb.org/mongo-driver/bson"
)

const databaseName = "variables"

const checkName = "easter_egg"

var variablesFilter = bson.M{
	"key": bson.M{"$nin": []string{"token", "password"}},
}

var eggFilter = bson.M{
	"key": checkName,
}

// Get variables of keys
func Get(keys ...string) (res Variables, err error) {
	res = make(Variables)

	data := make([]map[string]interface{}, 0)
	_, err = mongo.Find(
		"blotter",
		databaseName,
		bson.M{
			"key": bson.M{
				"$in": keys,
			},
		},
		nil,
		&data,
	)

	for _, d := range FromMapSliceToTypeSlice(data) {
		res[d.Key] = d.Value
	}
	return
}

// GetAll variables
func GetAll() (res BlotterVariables, err error) {
	defer errors.Wrapper(&err)

	temp := make([]map[string]interface{}, 0)
	if _, err = mongo.Find(
		"blotter",
		databaseName,
		variablesFilter,
		nil,
		&temp,
	); err != nil {
		return
	}
	res, err = NewBlotterVariables(temp)
	return
}

// SetMany variable
func SetMany(vars ...Type) (err error) {
	defer errors.Wrapper(&err)

	_, err = mongo.Remove(
		"blotter",
		databaseName,
		variablesFilter,
		nil,
	)
	if err != nil {
		return
	}

	_, err = mongo.Add(
		"blotter",
		databaseName,
		nil,
		transfer.ToInterfaceSlice(vars)...,
	)
	return
}

// Check easteregg and return the link if the key exists
func CheckEasterEgg(key string) (link string, miL int, maL int, err error) {
	var eggDict map[string]string
	temp := make([]map[string]interface{}, 0)

	if _, err = mongo.Find(
		"blotter",
		databaseName,
		eggFilter,
		nil,
		&temp,
	); err != nil {
		return
	}
	fmt.Println(key)
	// take the value of "easteregg" from the result of mongo
	oriEgg := temp[0]
	eggDict, miL, maL, err = SplitString(oriEgg["value"].(string))
	for k, v := range eggDict {
		tmpLength := len(k)
		// fmt.Println(k, tmpLength)
		if len(key) >= tmpLength && k == key[len(key)-tmpLength:] {
			link = v
		}
	}
	if _, ok := eggDict[key]; ok {
		link = eggDict[key]
	}
	return
}

// split the string with space and return the dict
// where the key is the first word and value is the second word
func SplitString(str string) (res map[string]string, miL int, maL int, err error) {
	miL = 100
	maL = 0
	res = make(map[string]string)
	str = strings.TrimSpace(str)
	stringlist := strings.Split(str, " ")
	max_length := len(stringlist)
	if max_length%2 != 0 {
		err = errors.New("The length of the string is wrong")
		return
	} else {
		for i := 0; i < max_length-1; i += 2 {
			res[stringlist[i]] = stringlist[i+1]
			miL = Min(miL, len(stringlist[i]))
			maL = Max(maL, len(stringlist[i]))
		}
	}
	return
}
func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
