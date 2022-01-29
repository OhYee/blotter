package email

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/OhYee/blotter/api/pkg/variable"
	"github.com/pkg/errors"
)

// Send email
func Send(hostWithPort, username, user, password string, ssl bool, subject, body string, to ...string) (err error) {
	defer func() {
		if err != nil {
			err = errors.Wrap(err, "Send email failed")
		}
	}()
	hp := strings.Split(hostWithPort, ":")
	host := hp[0]

	// make connection
	var conn net.Conn
	if ssl {
		conn, err = tls.Dial("tcp", hostWithPort, nil)
		if err != nil {
			return
		}
	} else {
		conn, err = net.Dial("tcp", hostWithPort)
		if err != nil {
			return
		}
	}

	// initiate connection
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return
	}

	// smtp auth
	auth := smtp.PlainAuth("", user, password, host)
	if err = client.Auth(auth); err != nil {
		return
	}

	// set sender
	if err = client.Mail(user); err != nil {
		return
	}

	// set receiver
	for _, receiver := range to {
		if err = client.Rcpt(receiver); err != nil {
			return
		}
	}

	// Write data
	w, err := client.Data()
	if err != nil {
		return
	}
	msg := []byte(fmt.Sprintf(
		"To: %s\r\nFrom: %s<%s>\r\nSubject: =?UTF-8?Q?%s?=\r\nContent-Type: text/html;charset=UTF-8\r\n\r\n%s",
		strings.Join(to, ","),
		username, user,
		subject, body,
	))
	if _, err = w.Write(msg); err != nil {
		return
	}
	if err = w.Close(); err != nil {
		return
	}

	err = client.Quit()
	return err
}

// GetSMTPData get smtp information from database
func GetSMTPData() (email, user, username, password, address string, ssl bool, root, blogName string, err error) {
	v, err := variable.Get(
		"email", "smtp_user", "smtp_password", "smtp_address", "smtp_ssl",
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
	if err = v.SetBool("smtp_ssl", &ssl, false); err != nil {
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
