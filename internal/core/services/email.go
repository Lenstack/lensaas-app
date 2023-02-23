package services

import (
	"bytes"
	"fmt"
	"gopkg.in/gomail.v2"
	"html/template"
	"strconv"
)

type IEmailService interface {
	Create(templateUrl string, to []string, subject string, body interface{}, attachments []string) (message *gomail.Message, err error)
	Send(mail *gomail.Message) error
}

type EmailService struct {
	Host     string
	Port     int
	Email    string
	Password string
}

func NewEmailService(host string, port string, email string, password string) *EmailService {
	portInt, _ := strconv.Atoi(port)
	return &EmailService{Host: host, Port: portInt, Email: email, Password: password}
}

func (es *EmailService) Create(templateUrl string, to []string, subject string, body interface{}, attachments []string) (message *gomail.Message, err error) {
	mail := gomail.NewMessage()
	mail.SetHeader("From", es.Email)
	mail.SetHeader("To", to...)
	mail.SetHeader("Subject", subject)

	if templateUrl != "" {
		var tpl bytes.Buffer

		emailTemplate, err := template.ParseFiles(templateUrl)
		if err != nil {
			return nil, err
		}

		err = emailTemplate.Execute(&tpl, body)
		if err != nil {
			return nil, err
		}
		mail.SetBody("text/html", tpl.String())
	} else {
		mail.SetBody("text/plain", body.(string))
	}

	for _, attachment := range attachments {
		mail.Attach(attachment)
	}
	return mail, nil
}

func (es *EmailService) Send(mail *gomail.Message) error {
	dialer := gomail.NewDialer(es.Host, es.Port, es.Email, es.Password)
	errChan := make(chan error)
	go func() {
		if err := dialer.DialAndSend(mail); err != nil {
			fmt.Println("Error sending email: ", err)
			errChan <- err
		}
	}()
	return nil
}
