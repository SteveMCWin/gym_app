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

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

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

		c.HTML(http.StatusOK, "profile.html", usr)

		// display user data somehow ig

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
		// this page is supposed to lead you to the HandlePostSignupSendMail()
		c.HTML(http.StatusOK, "signup.html", gin.H{})
	}
}

// func HandlePostSignup(db *models.DataBase) func(c *gin.Context) {
// 	return func(c *gin.Context) {
// 		email := c.PostForm("email")
// 		if email_exists := db.EmailExists(email); email_exists == false {
// 			log.Println("You already have an account!")
// 			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
// 			return
// 		}
//
// 		// sends email
// 		c.Redirect(http.StatusTemporaryRedirect, "/signup/send")
// 	}
// }

func HandlePostSignupSendMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		usr_email := c.PostForm("email")
		if email_exists := db.EmailExists(usr_email); email_exists == false {
			log.Println("You already have an account!")
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		token_val := CreateToken(usr_email, 3*time.Minute)

		new_mail := &mail.Mail{
			Recievers:    []string{usr_email},
			Subject:      "Signup Verification",
			TempaltePath: "./templates/test_mail.html",
			ExtLink:      "user/signup/from-mail/" + strconv.Itoa(token_val) + "/" + usr_email}

		err := mail.SendMailHtml(new_mail)
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}
		c.Redirect(http.StatusTemporaryRedirect, "user/signup/mail-sent")
	}
}

func HandleGetSignupMailSent() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render html that says the mail has been sent
		c.HTML(http.StatusOK, "mail_sent.html", gin.H{})
	}
}

func HandleGetSignupFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render forms with html and whatnot to get the user data such as usrname etc. etc.
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}
		usr_email := c.Param("email")

		if t_val, exists := signupTokens[token_val]; exists != true || t_val != usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		c.HTML(http.StatusOK, "account_creation.html", gin.H{"ID": token_val, "Email": usr_email})
	}
}

func HandlePostSignupFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get the form data which is the user's credentials and whatever and call db.create user and such
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}
		usr_email := c.Param("email")
		if t_val, exists := signupTokens[token_val]; exists == true && t_val == usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		training_since, err := time.Parse("2006-01-02", c.PostForm("training_since"))
		if err != nil {
			panic(err)
		}
		
		is_trainer := c.PostForm("is_trainer") != ""

		new_user := models.User{
			Id:            0,
			Name:          c.PostForm("name"),
			Email:         usr_email,
			Password:      c.PostForm("password"),
			TrainingSince: training_since,
			IsTrainer:     is_trainer,
			GymGoals:      c.PostForm("gym_goals"),
			CurrentGym:    c.PostForm("current_gym"),
		}

		usr_id, err := db.CreateUser(c, new_user)

		if err != nil {
			panic(err)
		}

		sessionManager.Put(c, "user_id", usr_id)

		delete(signupTokens, token_val)

		c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
		return
	}
}

// NOTE: Remember to limit the length of user password to less than 70 characters!

func HandleGetLogout(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionManager.Clear(c.Request.Context())

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}
