package utils

import (
	"bytes"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
)

type Email struct {
	Host     string
	Port     int
	Email    string
	Password string
}

func NewEmail(host string, port string, email string, password string) *Email {
	portInt, _ := strconv.Atoi(port)
	return &Email{Host: host, Port: portInt, Email: email, Password: password}
}

func (e *Email) Send(templateUrl string, to []string, subject string, body interface{}, attachments []string) error {
	mail := gomail.NewMessage()
	mail.SetHeader("From", e.Email)
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
		mail.SetBody("text/plain", body.(string))
	}

	for _, attachment := range attachments {
		mail.Attach(attachment)
	}

	dialer := gomail.NewDialer(e.Host, e.Port, e.Email, e.Password)

	errChan := make(chan error)
	go func() {
		if err := dialer.DialAndSend(mail); err != nil {
			errChan <- err
		}
	}()

	return nil
}
