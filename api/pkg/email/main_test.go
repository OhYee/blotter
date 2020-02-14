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
			email, user, username, password, address, root, blogName, err := GetSMTPData()
			if err != nil {
				t.Errorf("%s", err)
				t.FailNow()
			}

			err = Send(address, username, user, password, "测试", fmt.Sprintf(
				"测试邮件功能 %s <a href='%s'>%s</a>",
				time.Now().Format("2006-01-02 15:04:05"), root, blogName,
			), email)
			if err != nil {
				t.Errorf("%s", err)
				t.FailNow()
			}
			break
		}
	}
}
