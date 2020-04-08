package env

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/OhYee/rainbow/errors"
)

// GetEnv get environment from file
func GetEnv(filename string) (environments map[string]string, err error) {
	defer errors.Wrapper(&err)

	environments = make(map[string]string)

	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}

	for _, line := range strings.Split(string(b), "\n") {
		ss := strings.Split(line, "=")
		if len(ss) >= 2 {
			environments[ss[0]] = strings.Join(ss[1:], "=")
		}
	}
	return
}

// PWDFile file in current directory
func PWDFile(filename string) (absPath string) {
	var pwd string
	var err error
	if pwd, err = os.Getwd(); err != nil {
		pwd = ""
	}
	absPath = path.Join(pwd, filename)
	return
}
