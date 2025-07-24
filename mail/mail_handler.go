package mail

import (
	"bytes"
	// "github.com/joho/godotenv"
	"html/template"
	"log"
	"net/smtp"
	// "os"
)

var LoadedTempaltes map[string]*template.Template
var mail_sender string
var mail_pass string

type Mail struct { // NOTE: could rework this to have more stuff, such as the message text and whatnot so it's more reusable
	Recievers    []string
	Subject      string
	TempaltePath string
	ExtLink      string
}

func InitMail(m_pass, m_sender string) {
	mail_pass = m_pass
	mail_sender = m_sender

	LoadedTempaltes = make(map[string]*template.Template)
}

func SendMailHtml(mail *Mail) error {
	auth := smtp.PlainAuth(
		"",
		mail_sender,
		mail_pass,
		"smtp.gmail.com",
	)

	if _, exists := LoadedTempaltes[mail.TempaltePath]; exists == false {
		t, err := template.ParseFiles(mail.TempaltePath)
		if err != nil {
			log.Println("Fail 1")
			return err
		}

		LoadedTempaltes[mail.TempaltePath] = t

	}

	var body bytes.Buffer

	err := LoadedTempaltes[mail.TempaltePath].Execute(&body, mail)
	if err != nil {
		log.Println("Fail 2")
		return err
	}

	headers := "MIME-version: 1.0\r\nContent-Type: text/html; charset=\"UTF-8\"\r\n"

	msg := []byte((
		"Subject: " + mail.Subject + "\r\n" +
		headers + "\r\n" +
		body.String() + "\r\n"))

	err = smtp.SendMail(
		"smtp.gmail.com:587",
		auth,
		mail_sender,
		mail.Recievers,
		[]byte(msg),
	)

	return err
}
