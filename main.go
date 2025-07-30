package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"fitness_app/handlers"
	"fitness_app/mail"
	"fitness_app/models"
)

func main() {

	var db models.DataBase
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	err = db.CacheData()
	if err != nil {
		panic(err)
	}

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

	handler := handlers.SetUpRouter(domain, csrf_key, db)

	http.ListenAndServe(":8080", handler)
}
