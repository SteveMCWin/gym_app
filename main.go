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
	"github.com/gorilla/csrf"
	"github.com/joho/godotenv"
)

var sessionManager *scs.SessionManager

var domain string

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Couldn't load the .env")
	}

	domain = os.Getenv("DOMAIN")
	csrf_key := os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		log.Fatal("Couldn't load .env variables")
	}

	var db models.DataBase
	err = db.InitDatabase()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	sessionManager = scs.New()
	sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Store = sqlite3store.New(db.Data)
	sessionManager.Cookie.Persist = true
	sessionManager.Cookie.Secure = true

	router := gin.Default()

	router.LoadHTMLGlob("templates/*") // loads all templates from the templates directory

	router.GET("/", HandleGetHome())
	router.GET("/error-page", HandleGetError())

	user_router := router.Group("/user")

	user_router.GET("/profile", MiddlewareNoCache(), HandleGetProfile(&db))
	user_router.GET("/login", HandleGetLogin())
	user_router.POST("/login", HandlePostLogin(&db))
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup", HandleGetSignup())
	user_router.POST("/signup/send-mail", HandlePostSignupSendMail(&db))
	user_router.GET("/signup/mail-sent", HandleGetSignupMailSent())
	user_router.GET("/signup/from-mail/:id/:email", HandleGetSignupFromMail())
	user_router.POST("/signup/from-mail/:id/:email", HandlePostSignupFromMail(&db))
	user_router.GET("/delete_account", HandleGetDeleteAccount())
	user_router.POST("/delete_account", /* MiddlewareNoCache(), */ HandlePostDeleteAccount(&db))
	user_router.GET("/edit_profile", HandleGetEditProfile(&db))
	user_router.POST("/edit_profile", HandlePostEditProfile(&db))
	user_router.GET("/change_password", HandleGetChangePassword(&db))
	user_router.GET("/change_password/:id/:email", HandleGetChangePasswordFromMail())
	user_router.POST("/change_password/:id/:email", HandlePostChangePasswordFromMail(&db))

	handler := sessionManager.LoadAndSave(router)
	handler = csrf.Protect(
		[]byte(csrf_key),
		csrf.Secure(true),
	)(handler)

	http.ListenAndServe(":8080", handler)
}
