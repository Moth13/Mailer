package mailer

import (
	"net/smtp"

	"github.com/moth13/mailer/util"
)

type Mailer struct {
	Auth   smtp.Auth
	Config util.Config
}

type Email struct {
	To      string `json:"to" form:"to" binding:"required"`
	Subject string `json:"subject" form:"subject" binding:"required"`
	Body    string `json:"body" form:"body" binding:"required"`
}

func NewMailer(config util.Config) *Mailer {
	auth := smtp.PlainAuth("", config.From, config.Password, config.Host)
	return &Mailer{Auth: auth, Config: config}
}

func (m *Mailer) SendEmail(email Email) error {
	message := []byte("Subject: " + email.Subject + "\r\n\r\n" + email.Body)
	err := smtp.SendMail(m.Config.Host+":"+m.Config.Port, m.Auth, m.Config.From, []string{email.To}, message)
	if err != nil {
		return err
	}
	return nil
}
