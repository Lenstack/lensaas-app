package utils

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
)

type Email struct {
	host string
	port int
	user string
	pass string
}

func NewEmail(host string, port string, user string, pass string) *Email {
	portInt, _ := strconv.Atoi(port)
	return &Email{
		host: host,
		port: portInt,
		user: user,
		pass: pass,
	}
}

func (e *Email) Send(templateUrl string, from string, to []string, subject string, body, attachments []string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", from)
	mail.SetHeader("To", to...)
	mail.SetHeader("Subject", subject)

	if templateUrl != "" {
		var tpl bytes.Buffer

		emailTemplate, err := template.ParseFiles(templateUrl)
		if err != nil {
			return err
		}

		err = emailTemplate.Execute(&tpl, body)
		if err != nil {
			return err
		}
		mail.SetBody("text/html", tpl.String())
	} else {
		mail.SetBody("text/plain", body[0])
	}

	for _, attachment := range attachments {
		mail.Attach(attachment)
	}

	dialer := gomail.NewDialer(e.host, e.port, e.user, e.pass)

	go func() {
		if err := dialer.DialAndSend(mail); err != nil {
			return
		}
	}()

	return nil
}
