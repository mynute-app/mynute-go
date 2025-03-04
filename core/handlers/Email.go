package handlers

import (
	gomail "gopkg.in/mail.v2"
	"log"
)

type MailData struct {
	Username string
	Code     string
}

type MailType int

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}

type Configs struct {
	MailVerifTemplateID string
	PassResetTemplateID string
	SMTPHost            string
	SMTPPort            int
	SMTPUser            string
	SMTPPass            string
}

const (
	MailConfirmation = "MailConfirmation"
	PassReset        = "PassReset"
)

type MailService interface {
	SendMail(mailReq *Mail) error
	NewMail(from string, to []string, subject string, mailType MailType, data *MailData) *Mail
}

type GoMailService struct {
	log     *log.Logger
	configs *Configs
}

func (ms *GoMailService) SendEmail(mailReq *Mail) error {

	m := gomail.NewMessage()
	m.SetHeader("From", mailReq.from)

	m.SetHeader("To", mailReq.to...)

	// Set the subject and body based on the mail type
	// if mailReq.mtype == MailConfirmation {
	//     m.SetHeader("Subject", "Email Confirmation")
	//     m.SetBody("text/html", "Hello "+mailReq.data.Username+",<br><br>Your confirmation code is: "+mailReq.data.Code)
	// } else if mailReq.mtype == PassReset {
	//     m.SetHeader("Subject", "Password Reset")
	//     m.SetBody("text/html", "Hello "+mailReq.data.Username+",<br><br>Your password reset code is: "+mailReq.data.Code)
	// }

	// Set up the SMTP server
	d := gomail.NewDialer(ms.configs.SMTPHost, ms.configs.SMTPPort, ms.configs.SMTPUser, ms.configs.SMTPPass)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
