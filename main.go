package main

import (
	"net/http"

	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"fitness_app/models"
)

var sessionManager *scs.SessionManager

func main() {

	sessionManager = scs.New()
	// sessionManager.Lifetime = time.Hour * 24 * 30
	sessionManager.Lifetime = time.Minute * 2

	var db models.DataBase
	db.InitDatabase()

	router := gin.Default()
	user_router := router.Group("/user")
	
	user_router.GET("/profile")
	user_router.GET("/login")
	user_router.GET("/signup")

	handler := sessionManager.LoadAndSave(router)

	http.ListenAndServe(":8080", handler)
}
