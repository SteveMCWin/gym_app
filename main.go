package main

import (
	"log"
	"net/http"
	"os"

	"time"

	"fitness_app/models"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var sessionManager *scs.SessionManager

func main() {

	err := godotenv.Load()

	if err != nil {
		log.Fatal(err)
	}

	gmail_key := os.Getenv("GMAIL_APP_PASS")

	if gmail_key == "" {
		log.Fatal("No GMAIL_APP_PASS found")
	}

	var db models.DataBase
	err = db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionManager = scs.New()
	// sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Lifetime = time.Minute * 2   // NOTE: this one is just for debugging
	sessionManager.Store = sqlite3store.New(db.Data)

	router := gin.Default()

	router.GET("/")
	router.GET("/error-page")

	user_router := router.Group("/user")
	
	user_router.GET("/profile", HandleGetProfile(&db))
	user_router.GET("/login", HandleGetLogin())
	user_router.POST("/login", HandlePostLogin(&db))
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup")
	user_router.POST("/signup")
	user_router.GET("/signup/mail-sent")
	user_router.GET("/signup/from-mail/:id/:email")
	user_router.POST("/signup/from-mail")
	user_router.GET("/check_mail")

	handler := sessionManager.LoadAndSave(router)

	http.ListenAndServe(":8080", handler)
}
