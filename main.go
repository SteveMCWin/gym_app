package main

import (
	// "log"
	"net/http"
	// "os"

	"time"

	"fitness_app/models"

	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	// "github.com/joho/godotenv"
)

var sessionManager *scs.SessionManager

func main() {

	var db models.DataBase
	err := db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionManager = scs.New()
	// sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Lifetime = time.Minute * 2   // NOTE: this one is just for debugging
	sessionManager.Store = sqlite3store.New(db.Data)

	router := gin.Default()

	router.LoadHTMLGlob("templates/*")   // loads all templates from the templates directory

	router.GET("/", HandleGetHome())
	router.GET("/error-page")

	user_router := router.Group("/user")
	
	user_router.GET("/profile", HandleGetProfile(&db))
	user_router.GET("/login", HandleGetLogin())
	user_router.POST("/login", HandlePostLogin(&db))
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup", HandleGetSignup())
	user_router.POST("/signup/send-mail", HandlePostSignupSendMail(&db))
	user_router.GET("/signup/mail-sent", HandleGetSignupMailSent())
	user_router.GET("/signup/from-mail/:id/:email", HandleGetSignupFromMail())
	user_router.POST("/signup/from-mail/:id/:email", HandlePostSignupFromMail(&db))

	handler := sessionManager.LoadAndSave(router)

	http.ListenAndServe(":8080", handler)
}
