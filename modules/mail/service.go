package mail

import (
	"net/smtp"
	"strings"
)

type MailService interface {
	Send(to []string, cc []string, subject, message string) error
}

// impl `MailService` start

func (service *MailModule) Send(to []string, cc []string, subject, message string) error {
	body := "From: " + service.fromName + "\n" +
		"To: " + strings.Join(to, ",") + "\n" +
		"Cc: " + strings.Join(cc, ",") + "\n" +
		"Subject: " + subject + "\n\n" +
		message

	err := smtp.SendMail(service.address, service.auth, service.fromEmail, append(to, cc...), []byte(body))
	if err != nil {
		return err
	}

	return nil

}

// impl `MailService` end
