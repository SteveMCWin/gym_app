package main

import (
	"net/http"
	"log"
	"os"

	"github.com/joho/godotenv"

	"fitness_app/handlers"
	"fitness_app/models"
	"fitness_app/mail"
)

func main() {

	var db models.DataBase
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load the .env")
	}

	domain := os.Getenv("DOMAIN")
	csrf_key := os.Getenv("CSRF_KEY")

	mail_pass := os.Getenv("GMAIL_APP_PASS")
	mail_sender := os.Getenv("MAIL_SENDER")

	mail.InitMail(mail_pass, mail_sender)

	if domain == "" || csrf_key == "" {
		log.Fatal("Couldn't load .env variables")
	}

	router := handlers.SetUpRouter(domain, csrf_key, db)

	http.ListenAndServe(":8080", router)
}
