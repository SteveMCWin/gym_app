package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"fitness_app/mail"
	"fitness_app/models"

	"github.com/gin-gonic/gin"
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
			return
		}

		// display user data somehow ig
		_ = usr
	}
}

func HandleGetLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
			return
		}
		// just display the html ig using which you send a post request
	}
}

func HandlePostLogin(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		usr_id, err := db.AuthUser(email, password)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		sessionManager.Put(c, "user_id", usr_id)
		c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
	}
}

func HandleGetSignup() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
			return
		}
		// just display the html ig using which you send a post request
		// there is supposed to be an input field that takes in an email and a button for sending the user creation email
	}
}

func HandlePostSignupMailSent() func(c *gin.Context) {
	return func(c *gin.Context) {
		usr_email := c.PostForm("email")
		token_val := CreateToken(usr_email, 3 * time.Minute)
		new_mail := &mail.Mail{
			Recievers: []string{usr_email},
			Subject: "Signup Verification",
			TempaltePath: "./templates/test_mail.html",
			ExtLink: "user/signup/from-mail/"+strconv.Itoa(token_val)+"/"+usr_email}
		err := mail.SendMailHtml(new_mail)
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}
		c.Redirect(http.StatusTemporaryRedirect, "user/signup/mail-sent")
	}
}

func HandleGetSignupFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}
		usr_email := c.Param("email")

		if t_val, exists := signupTokens[token_val]; exists == true && t_val == usr_email {
			// render forms with html and whatnot to get the user data such as usrname etc. etc.
			delete(signupTokens, token_val)
			c.Redirect(http.StatusTemporaryRedirect, "user/signup/from-mail")
		}
	}
}

func HandlePostSignupFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get the form data which is the user's credentials and whatever and call db.create user and such
	}
}

func HandlePostSignup(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		if email_exists := db.EmailExists(email); email_exists == false {
			log.Println("You already have an account!")
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		// sends email
		c.Redirect(http.StatusTemporaryRedirect, "/user/check_mail")
	}
}

// NOTE: Remember to limit the length of user password to less than 70 characters!

func HandleGetLogout(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionManager.Clear(c.Request.Context())

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
