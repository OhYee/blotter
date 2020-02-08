package email

import (
	"fmt"
	"net/smtp"
	"strings"

	"github.com/OhYee/blotter/api/pkg/variable"
)

// Send email
func Send(host, username, user, password, subject, body string, to ...string) error {
	hp := strings.Split(host, ":")
	auth := smtp.PlainAuth("", user, password, hp[0])

	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s<%s>\r\nSubject: %s\r\nContent-Type: text/html;charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ","),
		username, user,

		subject, body,
	))
	err := smtp.SendMail(host, auth, user, to, msg)
	return err
}

// GetSMTPData get smtp information from database
func GetSMTPData() (email, user, username, password, address, root, blogName string, err error) {
	v, err := variable.Get(
		"email", "smtp_user", "smtp_password", "smtp_address",
		"smtp_username", "root", "blog_name",
	)
	if err != nil {
		return
	}

	if err = v.SetString("email", &email); err != nil {
		return
	}
	if err = v.SetString("smtp_user", &user); err != nil {
		return
	}
	if err = v.SetString("smtp_username", &username); err != nil {
		return
	}
	if err = v.SetString("smtp_password", &password); err != nil {
		return
	}
	if err = v.SetString("smtp_address", &address); err != nil {
		return
	}
	if err = v.SetString("root", &root); err != nil {
		return
	}
	if err = v.SetString("blog_name", &blogName); err != nil {
		return
	}
	return
}
