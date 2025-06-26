package mail

import (
	"bytes"
	"html/template"
	"net/smtp"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var LoadedTempaltes map[string]*template.Template
var mail_pass string
var mail_sender string

type Mail struct {
	Recievers []string
	Subject string
	TempaltePath string
	ExtLink string
}

func init() {
	LoadedTempaltes = make(map[string]*template.Template)

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	gmail_key := os.Getenv("GMAIL_APP_PASS")
	mail_sender := os.Getenv("MAIL_SENDER")

	if gmail_key == "" || mail_sender == ""{
		log.Fatal("ERROR: No mail sending data found")
	}
}

func SendMailHtml(mail *Mail) error {
	auth := smtp.PlainAuth(
		"",
		mail_sender,
		mail_pass,
		"stmp.gmail.com",
	)

	if _, exists := LoadedTempaltes[mail.TempaltePath]; exists == false {
		t, err := template.ParseFiles(mail.TempaltePath)
		if err != nil {
			return err
		}

		LoadedTempaltes[mail.TempaltePath] = t

	}

	var body bytes.Buffer

	err := LoadedTempaltes[mail.TempaltePath].Execute(&body, mail)
	if err != nil {
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
