package mail

import (
	"gopkg.in/gomail.v2"
)

type MailService interface {
	Send(to []string, cc []string, subject, message string) error
}

// impl `MailService` start

func (service *MailModule) Send(to []string, cc []string, subject, message string) error {
	mailer := gomail.NewMessage()

	mailer.SetHeader("From", service.fromEmail)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Cc", cc...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", message)

	err := service.dialer.DialAndSend(mailer)
	if err != nil {
		return err
	}

	return nil

}

// impl `MailService` end
