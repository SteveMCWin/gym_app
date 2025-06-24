package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"fitness_app/models"
)

func HandleGetProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
		}
		user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		usr, err := db.ReadUser(user_id)

		if err != nil {
			c.Redirect(http.StatusInternalServerError, "/error-page")
		}

		// display user data somehow ig
		_ = usr
	}
}

func HandleGetLogout(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionManager.Clear(c.Request.Context())

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
