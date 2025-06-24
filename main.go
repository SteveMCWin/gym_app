package main

import (
	"net/http"

	"time"

	"fitness_app/models"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/gin-gonic/gin"
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

	router.GET("/")
	router.GET("/error-page")

	user_router := router.Group("/user")
	
	user_router.GET("/profile", HandleGetProfile(&db))
	user_router.GET("/login")
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup")

	handler := sessionManager.LoadAndSave(router)

	http.ListenAndServe(":8080", handler)
}
