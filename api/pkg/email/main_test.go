package email

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestSend(t *testing.T) {
	for _, arg := range os.Args {
		if arg == "mail" {
			// email, user, username, password, address, root, blogName, err := GetSMTPData()
			email := "me@ohyee.cc"
			user := "sender@oyohyee.com"
			username := "sender@oyohyee.com"
			password := "QKpyTQc1RHJ8"
			address := "smtppro.zoho.com:465"
			root := "http://localhost:8080"
			blogName := "oyohyee"
			var err error = nil
			if err != nil {
				t.Errorf("%+v", err)
				t.FailNow()
			}

			err = Send(address, username, user, password, true, "测试", fmt.Sprintf(
				"测试邮件功能 %s <a href='%s'>%s</a>",
				time.Now().Format("2006-01-02 15:04:05"), root, blogName,
			), email)
			if err != nil {
				t.Errorf("%+v", err)
				t.FailNow()
			}
			break
		}
	}
}
