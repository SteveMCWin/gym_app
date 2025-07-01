package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"fitness_app/mail"
	"fitness_app/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
)

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	}
}

func HandleGetError() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.html", gin.H{})
	}
}

func HandleGetProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		// log.Println("user_id from sessionManager:", user_id)

		usr, err := db.ReadUser(user_id)

		if err != nil {
			log.Println("ERROR")
			log.Println(err)
			log.Println("ERROR")
			c.Redirect(http.StatusInternalServerError, "/error-page")
			return
		}

		user_view := gin.H {
			"Name": usr.Name,
			"Email": usr.Email,
			"TrainingSince": usr.TrainingSince.Format("2006-01-02"), // consider doing a .split on the string and rearrange
			"IsTrainer": usr.IsTrainer,
			"GymGoals": usr.GymGoals,
			"CurrentGym": usr.CurrentGym,
			"CurrentPlan": usr.CurrentPlan,
			"DateCreated": usr.DateCreated.Format("2006-01-02"), // consider doing a .split on the string and rearrange
		}

		c.HTML(http.StatusOK, "profile.html", user_view)

	}
}

func HandleGetLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
			return
		}

		c.HTML(http.StatusOK, "login.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request) })
	}
}

func HandlePostLogin(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		email := c.PostForm("email")
		password := c.PostForm("password")

		usr_id, err := db.AuthUserByEmail(email, password)

		if err != nil { // TODO: Handle errors better
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/user/login")
			return
		}

		sessionManager.Put(c.Request.Context(), "user_id", usr_id)
		c.Redirect(http.StatusSeeOther, "/user/profile")
	}
}

func HandleGetSignup() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
			return
		}

		c.HTML(http.StatusOK, "signup.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request) })
	}
}

func HandlePostSignupSendMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		usr_email := c.PostForm("email")
		if email_exists := db.EmailExists(usr_email); email_exists == true {
			log.Println("You already have an account!")
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		token_val := CreateToken(usr_email, 5*time.Minute)

		new_mail := &mail.Mail{
			Recievers:    []string{usr_email},
			Subject:      "Signup Verification",
			TempaltePath: "./templates/mail_register.html",
			ExtLink:      domain + "/user/signup/from-mail/" + strconv.Itoa(token_val) + "/" + usr_email} // NOTE: the domain mustn't end with a '/'

		err := mail.SendMailHtml(new_mail)
		if err != nil {
			// log.Fatalln(err) // WARNING: handle better than just panicing
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}
		c.Redirect(http.StatusSeeOther, "/user/signup/mail-sent")
	}
}

func HandleGetSignupMailSent() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render html that says the mail has been sent
		log.Println("IT GOT TO HandleGetSignupMailSent")
		c.HTML(http.StatusOK, "sent_mail.html", gin.H{})
	}
}

func HandleGetSignupFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render forms with html and whatnot to get the user data such as usrname etc. etc.
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		usr_email := c.Param("email")

		if t_val, exists := signupTokens[token_val]; exists != true || t_val != usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		c.HTML(http.StatusOK, "account_creation.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request), "ID": token_val, "Email": usr_email})
	}
}

func HandlePostSignupFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get the form data which is the user's credentials and whatever and call db.create user and such
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatalln(err) // WARNING: handle better than just panicing
		}

		log.Println("HandlePostSignupSendMail: got token_val")

		usr_email := c.Param("email")
		if t_val, exists := signupTokens[token_val]; exists != true || t_val != usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		log.Println("HandlePostSignupSendMail: got usr_email")

		training_since, err := time.Parse("2006-01-02", c.PostForm("training_since"))
		if err != nil {
			panic(err)
		}

		log.Println("HandlePostSignupSendMail: got parsed training_since")

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

		log.Println("HandlePostSignupSendMail: created user")

		sessionManager.Put(c.Request.Context(), "user_id", usr_id)

		log.Println("HandlePostSignupSendMail: created session cookie for user")

		delete(signupTokens, token_val)

		c.Redirect(http.StatusSeeOther, "/user/profile")
		return
	}
}

// NOTE: Remember to limit the length of user password to less than 70 characters! That is handled in the html and I hope that's enough

func HandleGetDeleteAccount() func(c *gin.Context) {
	return func(c *gin.Context) {
		var usr_id int

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			usr_id = 0
		} else {
			usr_id = sessionManager.GetInt(c.Request.Context(), "user_id")
		}

		c.HTML(http.StatusOK, "delete_accout.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request), "UserID": usr_id })
	}
}

func HandlePostDeleteAccount(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		password := c.PostForm("password")
		usr_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		log.Println("Got here")

		err := db.AuthUserByID(usr_id, password)
		if err != nil {
			log.Println("Couldn't delete accoutn")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/user/delete_account")
			return
		}

		db.DeleteUser(usr_id)
		sessionManager.Clear(c.Request.Context())
		log.Println("Deleted user with id", usr_id)
		c.Redirect(http.StatusSeeOther, "/")

	}
}

func MiddlewareNoCache() func(c *gin.Context) {
    return func(c *gin.Context) {
        c.Writer.Header().Set("Cache-Control", "no-store")
        c.Writer.Header().Set("Pragma", "no-cache")
        c.Writer.Header().Set("Expires", "0")
        c.Next()
    }
}

func HandleGetLogout(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionManager.Destroy(c.Request.Context())

		c.Redirect(http.StatusTemporaryRedirect, "/")
	}
}

func HandleGetEditProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		old_user, err := db.ReadUser(sessionManager.GetInt(c.Request.Context(), "user_id"))
		if err != nil {
			log.Println("COULDN'T GET USERS OLD DATA")
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
		}
		old_user_view := gin.H{
			csrf.TemplateTag: csrf.TemplateField(c.Request),
			"Name":          old_user.Name,
			"TrainingSince": old_user.TrainingSince.Format("2006-01-02"),
			"IsTrainer":     old_user.IsTrainer,
			"GymGoals":      old_user.GymGoals,
			"CurrentGym":    old_user.CurrentGym,
		}
		c.HTML(http.StatusOK, "edit_profile.html", old_user_view)
	}
}

func HandlePostEditProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		training_since, err := time.Parse("2006-01-02", c.PostForm("training_since"))
		if err != nil {
			panic(err)
		}

		log.Println("HandlePostSignupSendMail: got parsed training_since")

		is_trainer := c.PostForm("is_trainer") != ""

		edited_user := models.User{
			Id:            sessionManager.GetInt(c.Request.Context(), "user_id"),
			Name:          c.PostForm("name"),
			TrainingSince: training_since,
			IsTrainer:     is_trainer,
			GymGoals:      c.PostForm("gym_goals"),
			CurrentGym:    c.PostForm("current_gym"),
		}

		_, err = db.UpdateUserPublicData(&edited_user)

		if err != nil {
			log.Println("Couldn't edit user data?!")
		}

		c.Redirect(http.StatusSeeOther, "/user/profile")
	}
}

func HandleGetChangePassword(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		usr, err := db.ReadUser(sessionManager.GetInt(c.Request.Context(), "user_id"))

		if err != nil {
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		token_val := CreateToken(usr.Email, 5*time.Minute)

		new_mail := &mail.Mail{
			Recievers:    []string{usr.Email},
			Subject:      "Password Change",
			TempaltePath: "./templates/mail_change_password.html",
			ExtLink:      domain + "/user/change_password/" + strconv.Itoa(token_val) + "/" + usr.Email} // NOTE: the domain mustn't end with a '/'

		err = mail.SendMailHtml(new_mail)
		if err != nil {
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "sent_mail.html", gin.H{})
	}
}

func HandleGetChangePasswordFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		token_val, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}
		usr_email := c.Param("email")

		if t_val, exists := signupTokens[token_val]; exists != true || t_val != usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		c.HTML(http.StatusOK, "change_password.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request), "ID": token_val, "Email": usr_email })
	}
}

func HandlePostChangePasswordFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// NOTE: technically the cookie could expire while the user is using the website and reaching this point but scs should be handling that properly so it doesn't happen
		log.Println("Got to Password Change POST!!!!!")

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		new_password := c.PostForm("password")

		usr_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		_, err := db.UpdateUserPassword(usr_id, new_password)

		if err != nil {
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		log.Println("Updated user password successfully")

		token, err := strconv.Atoi(c.Param("id"))

		if err != nil {
			log.Println("ERROR")
			log.Println("COULDN'T DELETE SIGNUP TOKEN")
			log.Println("ERROR")
		} else {
			delete(signupTokens, token)
		}

		log.Println("Deleted token successfully")

		c.Redirect(http.StatusSeeOther, "/user/profile")
	}
}

func HandleGetCreatePlan() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
		}
		c.HTML(http.StatusOK, "make_plan.html", gin.H{ csrf.TemplateTag: csrf.TemplateField(c.Request) })
	}
}

func HandlePostCreatePlan(db *models.DataBase) func (c *gin.Context) {
	return func(c *gin.Context) {
		c.Request.ParseForm()
		form := c.Request.PostForm

		var plan_name string
		var plan_desc string
		// var set_to_current bool

		columns := make(map[int]string)
		cells := make(map[int][]string)

		fmt.Println()

		for key, vals := range form {
			value := vals[0]

			fmt.Println("Key:", key)
			fmt.Println("Vals:", vals)

			if strings.HasPrefix(key, "column_name_") {
				index, _ := strconv.Atoi(strings.TrimPrefix(key, "column_name_"))
				columns[index] = value
			} else if strings.HasPrefix(key, "cell_") {
				parts := strings.Split(strings.TrimPrefix(key, "cell_"), "_")
				if len(parts) != 2 {
					continue
				}
				colIdx, _ := strconv.Atoi(parts[0])
				rowIdx, _ := strconv.Atoi(parts[1])

				for len(cells[colIdx]) <= rowIdx {
					cells[colIdx] = append(cells[colIdx], "")
				}
				cells[colIdx][rowIdx] = value
			} else if key == "plan_name" {
				plan_name = value
			} else if key == "plan_description" {
				plan_desc = value
			} else if key == "make_current" {
				// set_to_current = value
			}
		}

		fmt.Println(columns)
		fmt.Println(cells)

		new_plan := models.WorkoutPlan {
			Name: plan_name,
			Description: plan_desc,
			Creator: sessionManager.GetInt(c.Request.Context(), "user_id"),
		}

		wp_id, err := db.CreateWorkoutPlan(c, &new_plan)

		if err != nil {
			c.Redirect(http.StatusSeeOther, "/error-page")
		}

		for i, day_name := range columns {
			for j, ex_name := range cells[i] {
				// Gotta figure out how Imma get the exercise id, other than that we shoul be ok
				_ = day_name
				_ = j
				_ = ex_name
				_ = wp_id
			}
		}


		// Use columns + cells as needed
		c.String(http.StatusOK, "Table submitted with %d columns", len(columns))
	}
}
