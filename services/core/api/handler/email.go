package handler

import (
	"log"
	"os"
	"strconv"

	gomail "gopkg.in/mail.v2"
)

type MailData struct {
	Clientname string
	Code       string
}

type Mail struct {
	from    string
	to      []string
	subject string
	body    string
	mtype   MailType
	data    *MailData
}
type MailType int

type Configs struct {
	MailVerifTemplateID string
	PassResetTemplateID string
	SMTPHost            string
	SMTPPort            int
	SMTPClient          string
	SMTPPass            string
}

const (
	MailConfirmation = "MailConfirmation"
	PassReset        = "PassReset"
)

var EmailRequestTest = &Mail{
	from:    "hello@demomailtrap.com",
	to:      []string{"luigiazoreng@gmail.com"},
	subject: "Test",
	body:    "Test",
	mtype:   2,
	data: &MailData{Clientname: "Test Client",
		Code: "123456"},
}

type MailService interface {
	SendMail(mailReq *Mail) error
	NewMail(from string, to []string, subject string, mailType int, data *MailData) *Mail
}

type GoMailService struct {
	configs *Configs
}

func NewGoMailService(configs *Configs) *GoMailService {
	setConfigs(configs)
	return &GoMailService{configs}

}
func (ms *GoMailService) SendEmail(mailReq *Mail) error {

	message := gomail.NewMessage()
	message.SetHeader("From", mailReq.from)
	message.SetHeader("To", mailReq.to...)

	log.Println("Sending email to: ", mailReq.to)
	message.SetBody("text/html", mailReq.body)
	message.SetHeader("Subject", mailReq.subject)
	// Set the subject and body based on the mail type
	// if mailReq.mtype == MailConfirmation {
	//     m.SetHeader("Subject", "Email Confirmation")
	//     m.SetBody("text/html", "Hello "+mailReq.data.Clientname+",<br><br>Your confirmation code is: "+mailReq.data.Code)
	// } else if mailReq.mtype == PassReset {
	//     m.SetHeader("Subject", "Password Reset")
	//     m.SetBody("text/html", "Hello "+mailReq.data.Clientname+",<br><br>Your password reset code is: "+mailReq.data.Code)
	// }

	// Set up the SMTP server
	d := gomail.NewDialer(ms.configs.SMTPHost, ms.configs.SMTPPort, ms.configs.SMTPClient, ms.configs.SMTPPass)

	// Send the email
	if err := d.DialAndSend(message); err != nil {
		return err
	}

	return nil
}

func setConfigs(configs *Configs) {
	var err error
	configs.SMTPPort, err = strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		log.Fatal("SMTP_PORT must be an integer")
	}
	configs.SMTPHost = os.Getenv("SMTP_HOST")
	configs.SMTPClient = os.Getenv("SMTP_USER")
	configs.SMTPPass = os.Getenv("SMTP_PASS")
}

