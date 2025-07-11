package mailer

import (
	"net/smtp"
)

type Mailer struct {
	smtpHost string
	smtpPort string
	from     string
	password string
}

func NewMailer(smtpHost, smtpPort, from, password string) *Mailer {
	return &Mailer{
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		from:     from,
		password: password,
	}
}

func (m *Mailer) SendEmail(email, message string) error {
	to := []string{
		email,
	}

	auth := smtp.PlainAuth("", m.from, m.password, m.smtpHost)

	err := smtp.SendMail(m.smtpHost+":"+m.smtpPort, auth, m.from, to, []byte(message))
	if err != nil {
		return err
	}
	return nil
}
