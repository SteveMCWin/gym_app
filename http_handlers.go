package main

import (
	"log"
	"net/http"
	"strconv"
	"time"
	// "slices"

	"fitness_app/mail"
	"fitness_app/models"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/csrf"
	// "github.com/lithammer/fuzzysearch/fuzzy"
)

func HandleGetHome() func(c *gin.Context) {
	return func(c *gin.Context) {
		var user_id int
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			user_id = -1
		} else {
			user_id = sessionManager.GetInt(c.Request.Context(), "user_id")
		}
		c.HTML(http.StatusOK, "index.html", gin.H{"user_id": user_id})
	}
}

func HandleGetError() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "error.html", gin.H{})
	}
}

func HandleGetPing() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	}
}

func HandleGetCurrentProfile() func(c *gin.Context) { // this is used for redirecting to the users own profile
	return func(c *gin.Context) {

		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!")
		log.Println("CAME TO USER PROFILE FIRST")
		log.Println("!!!!!!!!!!!!!!!!!!!!!!!!!!")

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		c.Redirect(http.StatusTemporaryRedirect, "/user/"+strconv.Itoa(user_id))
	}
}

func HandleGetProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusInternalServerError, "/error-page")
			return
		}

		if user_id == 0 {
			user_id = requesting_user_id
		}

		usr, err := db.ReadUser(user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusInternalServerError, "/error-page")
			return
		}

		current_plan, err := db.ReadWorkoutPlan(usr.CurrentPlan)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "profile.html", gin.H{ "usr": usr, "requesting_user_id": requesting_user_id, "current_plan": current_plan })
	}
}

func HandleGetLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/"+strconv.Itoa(sessionManager.GetInt(c.Request.Context(), "user_id")))
			return
		}

		c.HTML(http.StatusOK, "login.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request)})
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
		c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(usr_id))
	}
}

func HandleGetSignup() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") {
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
			return
		}

		c.HTML(http.StatusOK, "signup.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request)})
	}
}

func HandlePostSignup(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		usr_email := c.PostForm("email")
		log.Println("User's email:", usr_email)
		if email_exists := db.EmailExists(usr_email); email_exists == true {
			log.Println("You already have an account!")
			c.Redirect(http.StatusSeeOther, "/user/login")
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
			log.Println("FAILLLLED TO SEND MAILLLL")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}
		c.Redirect(http.StatusSeeOther, "/mail-sent")
	}
}

func HandleGetSignupMailSent() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render html that says the mail has been sent
		c.HTML(http.StatusOK, "sent_mail.html", gin.H{})
	}
}

func HandleGetSignupFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		// render forms with html and whatnot to get the user data such as usrname etc. etc.
		token_val, err := strconv.Atoi(c.Param("token_id"))
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

		c.HTML(http.StatusOK, "account_creation.html", gin.H{
			csrf.TemplateTag: csrf.TemplateField(c.Request),
			"ID": token_val,
			"Email": usr_email,
			"Gyms": models.FetchAllCachedGyms(),
		})
	}
}

func HandlePostSignupFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// get the form data which is the user's credentials and whatever and call db.create user and such
		token_val, err := strconv.Atoi(c.Param("token_id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		usr_email := c.Param("email")
		if t_val, exists := signupTokens[token_val]; exists != true || t_val != usr_email {
			log.Println("ERROR: invalid token or token value")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		training_since, err := time.Parse("2006-01-02", c.PostForm("training_since"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		is_trainer := c.PostForm("is_trainer") != ""

		curr_gym_id, err := strconv.Atoi(c.PostForm("current_gym"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		curr_gym, ok := models.FetchCachedGym(curr_gym_id)
		if !ok {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		new_user := models.User{
			Id:            0,
			Name:          c.PostForm("name"),
			Email:         usr_email,
			Password:      c.PostForm("password"),
			TrainingSince: training_since,
			IsTrainer:     is_trainer,
			GymGoals:      c.PostForm("gym_goals"),
			CurrentGym:    *curr_gym,
		}

		usr_id, err := db.CreateUser(c, new_user)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		sessionManager.Put(c.Request.Context(), "user_id", usr_id)

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

		c.HTML(http.StatusOK, "delete_accout.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request), "UserID": usr_id})
	}
}

func HandlePostDeleteAccount(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		password := c.PostForm("password")
		usr_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		err := db.AuthUserByID(usr_id, password)
		if err != nil {
			log.Println("Couldn't delete accoutn")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/user/delete_account")
			return
		}

		db.DeleteUser(usr_id)
		sessionManager.Clear(c.Request.Context())
		// log.Println("Deleted user with id", usr_id)
		c.Redirect(http.StatusSeeOther, "/")

	}
}

func HandleGetLogout(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		sessionManager.Destroy(c.Request.Context())

		c.Redirect(http.StatusTemporaryRedirect, "/user/login")
	}
}

func HandleGetEditProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		old_user, err := db.ReadUser(sessionManager.GetInt(c.Request.Context(), "user_id"))
		if err != nil {
			log.Println("COULDN'T GET USERS OLD DATA")
			c.Redirect(http.StatusTemporaryRedirect, "/user/profile")
		}

		c.HTML(http.StatusOK, "edit_profile.html", gin.H{
			csrf.TemplateTag: csrf.TemplateField(c.Request),
			"old_user": old_user,
			"Gyms": models.FetchAllCachedGyms(), // NOTE: update the html side too
		})
	}
}

func HandlePostEditProfile(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {


		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		// user_id_param := c.Param("id")
		// user_id, err := strconv.Atoi(user_id_param)
		// if err != nil {
		// 	log.Println(err)
		// 	c.Redirect(http.StatusInternalServerError, "/error-page")
		// 	return
		// }
		//
		// if requesting_user_id != user_id {
		// 	log.Println("You cannot edit other peoples profiles")
		// 	c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(requesting_user_id))
		// }

		training_since, err := time.Parse("2006-01-02", c.PostForm("training_since"))
		if err != nil {
			panic(err)
		}

		is_trainer := c.PostForm("is_trainer") != ""

		curr_gym_id, err := strconv.Atoi(c.PostForm("current_gym"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		curr_gym, ok := models.FetchCachedGym(curr_gym_id)
		if !ok {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		edited_user := models.User{
			Id:            requesting_user_id,
			Name:          c.PostForm("name"),
			TrainingSince: training_since,
			IsTrainer:     is_trainer,
			GymGoals:      c.PostForm("gym_goals"),
			CurrentGym:    *curr_gym,
		}

		_, err = db.UpdateUserPublicData(&edited_user)

		if err != nil {
			log.Println("Couldn't edit user data?!")
		}

		c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(requesting_user_id))
	}
}

func HandleGetChangePassword(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			log.Println("No session token found")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		usr, err := db.ReadUser(sessionManager.GetInt(c.Request.Context(), "user_id"))

		if err != nil {
			log.Println("Couldn't find user in database")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		token_val := CreateToken(usr.Email, 5*time.Minute)

		log.Println("Sending change pass mail to", usr.Email)

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

		log.Println("SENT CHANGE PASSWORD MAIL")

		c.HTML(http.StatusOK, "sent_mail.html", gin.H{})
	}
}

func HandleGetChangePasswordFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		token_val, err := strconv.Atoi(c.Param("token_id"))
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

		c.HTML(http.StatusOK, "change_password.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request), "ID": token_val, "Email": usr_email})
	}
}

func HandlePostChangePasswordFromMail(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// if sessionManager.Exists(c.Request.Context(), "user_id") == false {
		// 	c.Redirect(http.StatusSeeOther, "/error-page")
		// 	return
		// }

		new_password := c.PostForm("password")

		usr_email := c.Param("email")

		usr_id, err := db.ReadUserIdByEmail(usr_email)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		// usr_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		_, err = db.UpdateUserPassword(usr_id, new_password)

		if err != nil {
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		log.Println("Updated user password successfully")

		token, err := strconv.Atoi(c.Param("token_id"))

		if err != nil {
			log.Println("ERROR")
			log.Println("COULDN'T DELETE SIGNUP TOKEN")
			log.Println("ERROR")
		} else {
			delete(signupTokens, token)
		}

		log.Println("Deleted token successfully")

		c.Redirect(http.StatusSeeOther, "/user/logout")
	}
}

func HandleGetCreatePlan() func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		all_exercises := models.FetchAllCachedExercises()

		c.HTML(http.StatusOK, "make_plan.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request), "all_exercises": all_exercises})
	}
}

func HandlePostCreatePlan(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		usr_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		var plan models.WorkoutPlan
		if err := c.BindJSON(&plan); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		plan.Creator = usr_id

		wp_id, err := db.CreateWorkoutPlan(&plan)

		if err != nil {
			log.Println("ERROR")
			log.Println(err)
			log.Println("ERROR")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if plan.MakeCurrent {
			_, err = db.UpdateUserCurrentPlan(usr_id, wp_id)
		}

		if err != nil {
			log.Println("ERROR")
			log.Println(err)
			log.Println("ERROR")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}


		c.JSON(http.StatusOK, wp_id)
		// c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(usr_id))
	}
}

func HandleGetViewCurrentPlan(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			log.Println("Must be logged in to view current plan")
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		user_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		if user_id <= 0 {
			user_id = sessionManager.GetInt(c.Request.Context(), "user_id")
		}

		user, err := db.ReadUser(user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/user/"+strconv.Itoa(user_id)+"/plan/view/"+strconv.Itoa(user.CurrentPlan))
	}
}

func HandleGetViewPlan(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}
		if user_id <= 0 {
			user_id = requesting_user_id
		}

		user, err := db.ReadUser(user_id) // needed for checking if the plan being viewed is the users current plan
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		wp_id, err := strconv.Atoi(c.Param("wp_id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		wp, err := db.ReadWorkoutPlan(wp_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}
		if wp.Id == 1 {
			log.Println("Cannot view placeholder plan")
			c.Redirect(http.StatusTemporaryRedirect, "/user/0/plan/create")
			return
		}

		makeCurrent := user.CurrentPlan != wp_id && requesting_user_id == user_id

		plan_analysis := wp.GetAnalysis()

		c.HTML(http.StatusOK, "view_plan.html", gin.H{
			"wp": wp,
			"MakeCurrent": makeCurrent,
			"PlanAnalysis": plan_analysis,
			"user_id": user_id,
			"requesting_user_id": requesting_user_id,
		}) // WARNING: consider adding csrf protection especially if you enable editing the plan
	}
}

func HandleGetViewAllUserPlans(Db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		user_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		if user_id <= 0 {
			user_id = sessionManager.GetInt(c.Request.Context(), "user_id")
		}

		wps, err := Db.ReadAllWorkoutsUserUses(user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "view_all_user_plans.html", gin.H{ "wps": wps, "user_id": user_id })
	}
}

func HandleGetMakePlanCurrent(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		if requesting_user_id != user_id {
			log.Println("Cannot set someone else's plan as your current plan!")
			c.Redirect(http.StatusTemporaryRedirect, "/")
			return
		}

		wp_id, err := strconv.Atoi(c.Param("wp_id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		_, err = db.UpdateUserCurrentPlan(requesting_user_id, wp_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/user/"+strconv.Itoa(requesting_user_id))
			return
		}

		c.Redirect(http.StatusTemporaryRedirect, "/user/"+strconv.Itoa(requesting_user_id))
	}

}

func HandleGetEditPlan(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		user_id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		if user_id != sessionManager.GetInt(c.Request.Context(), "user_id") {
			log.Println("Cannot edit a plan that isn't yours!!")
			c.Redirect(http.StatusTemporaryRedirect, "/")
		}

		wp_id, err := strconv.Atoi(c.Param("wp_id"))
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		wp, err := db.ReadWorkoutPlan(wp_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}
		if wp.Id == 1 {
			log.Println("Cannot view placeholder plan")
			c.Redirect(http.StatusTemporaryRedirect, "/user/0/create")
			return
		}

		c.HTML(http.StatusOK, "edit_plan.html", gin.H{
			csrf.TemplateTag: csrf.TemplateField(c.Request),
			"wp": wp,
			"all_exercises": models.FetchAllCachedExercises(),
			"user_id": user_id,
		})
	}
}

func HandlePostEditPlan(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		user_id_param := c.Param("id")
		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		if user_id != sessionManager.GetInt(c.Request.Context(), "user_id") {
			log.Println("Cannot edit a plan that isn't yours!!")
			c.Redirect(http.StatusSeeOther, "/")
		}

		wp_id_param := c.Param("wp_id")

		var edited_wp models.WorkoutPlan
		if err := c.BindJSON(&edited_wp); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		edited_wp.Id, err = strconv.Atoi(wp_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		_, err = db.UpdateWorkoutPlan(&edited_wp)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.Redirect(http.StatusSeeOther, "/user/"+user_id_param+"/plan/view/"+wp_id_param)
	}
}

func HandleGetForgotPassword() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "forgot_password.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request)}) // may not need csrf protection here
	}
}

func HandlePostForgotPassword(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		usr_email := c.PostForm("email")
		log.Println("User's email:", usr_email)
		if email_exists := db.EmailExists(usr_email); email_exists != true {
			log.Println("There is no user with the email", usr_email)
			c.Redirect(http.StatusSeeOther, "/user/forgot_password")
			return
		}

		token_val := CreateToken(usr_email, 5*time.Minute)

		new_mail := &mail.Mail{
			Recievers:    []string{usr_email},
			Subject:      "Password Change",
			TempaltePath: "./templates/mail_change_password.html",
			ExtLink:      domain + "/user/forgot_password/from-mail/" + strconv.Itoa(token_val) + "/" + usr_email} // NOTE: the domain mustn't end with a '/'

		err := mail.SendMailHtml(new_mail)
		if err != nil {
			log.Println("FAILLLLED TO SEND MAILLLL")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}
		c.Redirect(http.StatusSeeOther, "/mail-sent")
	}
}

func HandleGetChangePassFromMail() func(c *gin.Context) {
	return func(c *gin.Context) {
		token_val, err := strconv.Atoi(c.Param("token_id"))
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

		c.HTML(http.StatusOK, "change_password.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request), "ID": token_val, "Email": usr_email})
	}
}

func HandleGetTracks(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("user_id")
		var user_id int
		var err error
		if user_id_param == "" || user_id_param == "0"{
			user_id = requesting_user_id
		} else {
			user_id, err = strconv.Atoi(c.Param("user_id"))
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusSeeOther, "/error-page")
				return
			}
		}

		workout_tracks, err := db.ReadUsersWorkoutTracks(user_id, requesting_user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "users_tracks.html", gin.H{
			"tracks": workout_tracks,
			"requesting_user_id": user_id,
		})
	}
}

func HandleGetTracksCreate(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}

		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		plans, err := db.ReadUsersRecentlyTrackedPlans(requesting_user_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "create_track.html", gin.H{
			csrf.TemplateTag: csrf.TemplateField(c.Request),
			"plans": plans,
			"user_id": requesting_user_id,
		})
	}
}

func HandlePostTracksCreate(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		p_id := c.Param("plan_id")
		if p_id == "" {
			log.Println("No plan id??")
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		plan_id, err := strconv.Atoi(p_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		make_private := c.PostForm("make_private_"+p_id) != ""

		p, ok := models.FetchCachedPlanBasic(plan_id)
		if !ok {
			log.Println("Couldn't fetch cached plan of id:", plan_id)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}
		wt := models.WorkoutTrack{
			Plan:        *p,
			User:        user_id,
			IsPrivate:   make_private,
			WorkoutDate: time.Now(),
		}

		wt.Id, err = db.CreateWorkoutTrack(&wt)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		err = db.CreateTrackDataForTrack(&wt)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(user_id)+"/track/view/"+strconv.Itoa(wt.Id))
	}
}

func HandleGetViewTrack(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("id")
		wt_id_param := c.Param("track_id")

		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}
		wt_id, err := strconv.Atoi(wt_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		track, err := db.ReadWorkoutTrack(wt_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if track.IsPrivate && user_id != requesting_user_id {
			log.Println("NUH UUUUH")
			c.Redirect(http.StatusTemporaryRedirect, "/error-page") // NOTE: create a page for Private or something
			return
		}

		track_data, err := db.ReadTrackDataForTrack(track.Id)
		if err != nil {
			log.Println("ERROR: error while reading track data")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "view_track.html", gin.H{"track_data": track_data, "Days": track.ExDays, "user_id": user_id, "requesting_user_id": requesting_user_id})
	}
}

func HandleGetTracksEdit(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("id")
		wt_id_param := c.Param("track_id")

		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println("Get 1")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if requesting_user_id != user_id {
			log.Println("NU UUUUH")
			c.Redirect(http.StatusTemporaryRedirect, "/error-page") // NOTE: create a page for Private or something
			return
		}

		wt_id, err := strconv.Atoi(wt_id_param)
		if err != nil {
			log.Println("Get 2")
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		wt, err := db.ReadWorkoutTrack(wt_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if wt.User != user_id {
			log.Println("NU UUUUH")
			c.Redirect(http.StatusTemporaryRedirect, "/error-page") // NOTE: create a page for Private or something
			return
		}

		td, err := db.ReadTrackDataForTrack(wt_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.HTML(http.StatusOK, "edit_track.html", gin.H{csrf.TemplateTag: csrf.TemplateField(c.Request), "Days": wt.ExDays, "track_data": td, "user_id": user_id})
	}
}

func HandlePostTracksEdit(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("id")
		wt_id_param := c.Param("track_id")

		user_id, err := strconv.Atoi(user_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if requesting_user_id != user_id {
			log.Println("NU UUUUH")
			c.Redirect(http.StatusTemporaryRedirect, "/error-page")
			return
		}

		wt_id, err := strconv.Atoi(wt_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		wt, err := db.ReadWorkoutTrack(wt_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		if wt.User != user_id {
			log.Println("NU UUUUH")
			c.Redirect(http.StatusTemporaryRedirect, "/error-page") // NOTE: create a page for Private or something
			return
		}

		var track_json models.TrackJSON
		if err := c.BindJSON(&track_json); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		_, err = db.UpdateMultipleTrackData(track_json.Data)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusSeeOther, "/error-page")
			return
		}

		c.Redirect(http.StatusSeeOther, "/user/"+strconv.Itoa(user_id)+"/track/view/"+wt_id_param)
	}
}

func HandleGetTracksViewLatest(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {


	}
}

func HandleGetTracksDelete(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		if sessionManager.Exists(c.Request.Context(), "user_id") == false {
			c.Redirect(http.StatusTemporaryRedirect, "/user/login")
			return
		}
		requesting_user_id := sessionManager.GetInt(c.Request.Context(), "user_id")

		user_id_param := c.Param("id")
		// user_id, err := strconv.Atoi(user_id_param)
		// if err != nil {
		// 	log.Println(err)
		// 	c.Redirect(http.StatusSeeOther, "/error-page")
		// 	return
		// }

		track_id_param := c.Param("track_id")
		track_id, err := strconv.Atoi(track_id_param)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		track, err := db.ReadWorkoutTrack(track_id)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		if requesting_user_id != track.User /* || user_id != track.User */{
			log.Println("You cannot delete another persons track")
		}

		_, err = db.DeleteWorkoutTrack(track)
		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		c.Redirect(http.StatusPermanentRedirect, "/user/"+user_id_param+"/track/view_all")
	}
}

func HandleGetSearchForUser(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		// just render some html for the search page
		query := c.Query("name")
		if query == "" {
			// render page
			c.HTML(http.StatusOK, "search_users.html", gin.H{})
		} else {
			// return JSON results
			results, err := db.SearchForUsers(query, sessionManager.GetInt(c.Request.Context(), "user_id"))
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusPermanentRedirect, "/error-page")
				return
			}
			c.JSON(200, results)
		}
	}
}

func HandleGetPlanJSON(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {
		wp_id_param := c.Param("wp_id")
		wp_id, err := strconv.Atoi(wp_id_param)
		if err != nil {
			log.Println("ERROR")
			log.Println(err)
			log.Println("ERROR")
			c.JSON(http.StatusInternalServerError, gin.H{})
		}

		wp, err := db.ReadWorkoutPlan(wp_id)
		if err != nil {
			log.Println("ERROR")
			log.Println(err)
			log.Println("ERROR")
			c.JSON(http.StatusInternalServerError, gin.H{})
		}

		// WARNING: Really bad performance-vise. Try a solution with a custom MarshalJSON or something. But for now it works
		for _, d := range wp.Days {
			for _, e := range d.Exercises {
				for _, t := range e.Exercise.Targets {
					t.Exercises = nil // Needed, otherwise the json gets into an infinite cycle because exercise has []Target and target has []Exercise
				}
			}
		}

		c.JSON(http.StatusOK, wp)
	}
}

func HandleGetViewAllGyms() func(c *gin.Context) {
	return func(c *gin.Context) {
		gyms := models.FetchAllCachedGyms()

		c.HTML(http.StatusOK, "view_all_gyms.html", gin.H{ "gyms": gyms })
	}
}

func HandleGetViewGym(db *models.DataBase) func(c *gin.Context) {
	return func(c *gin.Context) {

		gym_id_param := c.Param("gym_id")
		gym_id, err := strconv.Atoi(gym_id_param)

		if err != nil {
			log.Println(err)
			c.Redirect(http.StatusPermanentRedirect, "/error-page")
			return
		}

		var user_has_plan bool
		var ex_no_eq []int

		if sessionManager.Exists(c.Request.Context(), "user_id") {
			user_id := sessionManager.GetInt(c.Request.Context(), "user_id")
			user, err := db.ReadUser(user_id)
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusPermanentRedirect, "/error-page")
				return
			}

			user_has_plan = true

			ex_no_eq, _, err = db.CheckIfGymHasPlanEquipment(gym_id, user.CurrentPlan)
			if err != nil {
				log.Println(err)
				c.Redirect(http.StatusPermanentRedirect, "/error-page")
				return
			}
		}

		c.HTML(http.StatusOK, "view_gym.html", gin.H {
			"user_has_plan": user_has_plan,
			"ex_no_eq": ex_no_eq,
		})
	}
}

// MIDDLEWARE

func MiddlewareNoCache() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-store")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	}
}
