package mailer

import (
	"net/smtp"

	"github.com/moth13/mailer/models"
	"github.com/moth13/mailer/util"
)

type Mailer struct {
	Auth   smtp.Auth
	Config util.Config
}

func NewMailer(config util.Config) *Mailer {
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	return &Mailer{Auth: auth, Config: config}
}

func (m *Mailer) SendEmail(email models.Email) error {
	message := []byte("Subject: " + email.Subject + "\r\n\r\n" + email.Body)
	err := smtp.SendMail(m.Config.Host+":"+m.Config.Port, m.Auth, m.Config.From, []string{email.To}, message)
	if err != nil {
		email.Status = models.StatusErrorAtSend
		return err
	}
	email.Status = models.StatusSent
	return nil
}
