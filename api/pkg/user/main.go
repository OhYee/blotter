package user

import (
	"crypto/sha512"
	"encoding/hex"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/OhYee/blotter/mongo"
	"github.com/OhYee/blotter/output"

	"go.mongodb.org/mongo-driver/bson"
)

// PasswordHash get the hash of the password
func PasswordHash(username string, password string, userID string) (h string) {
	hash := sha512.New()

	// Can not use username as hash salt
	runes := []rune(password)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	reversePassword := string(runes)

	hash.Write([]byte(password + userID + reversePassword))
	h = hex.EncodeToString(hash.Sum(nil))
	output.Debug("%+v", h)
	return
}

// CheckToken check the token is valid
func CheckToken(token string) bool {
	m, err := variable.Get("token")
	if err != nil {
		return false
	}
	_token := ""
	m.SetString("token", &_token)
	if _token != token {
		return false
	}
	return true
}

// DeleteToken delete the token
func DeleteToken() {
	mongo.Remove("blotter", "variables", bson.M{"key": "token"}, nil)
}
