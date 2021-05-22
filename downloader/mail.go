package downloader

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type MailConf struct {
	FromEmail string
	FromName string
	ToEmail string
	Subject string
	Message string
	Attachment string
	Password  string
	SmtpHost string
	SmtpPort int
	AdminEmail string
}

type Mail struct {
	MailConf
}

func NewMail(conf MailConf)  *Mail {
	return &Mail{conf}
}

func (mail *Mail)SendMail() {
	m := gomail.NewMessage()
	m.SetHeader("From", fmt.Sprintf("%s <%s>", mail.FromName, mail.FromEmail))
	m.SetHeader("To", mail.ToEmail)
	m.SetHeader("Subject", mail.Subject)
	m.SetBody("text/html", mail.Message)
	m.Attach(mail.Attachment)

	d := gomail.NewDialer(mail.SmtpHost, mail.SmtpPort, mail.AdminEmail, mail.Password)
	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}
}
