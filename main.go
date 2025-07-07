package main

import (
	"log"
	"net/http"
	"html/template"
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

// used by gin to load template funcs
func templateFuncs() template.FuncMap {
    return template.FuncMap{
        "until": func(n int) []int {
			result := make([]int, n)
			for i := range n {
				result[i] = i
			}
			return result
        },
    }
}

// TODO: refactor the golang structs
// TODO: caching
// TODO: actual exercise name in track view

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

    router.SetFuncMap(templateFuncs())
	router.LoadHTMLGlob("templates/*") // loads all templates from the templates directory

	router.GET("/", HandleGetHome())
	router.GET("/error-page", HandleGetError())
	router.GET("/mail-sent", HandleGetSignupMailSent())

	user_router := router.Group("/user")
	// profile_router := router.Group("/user/profile") // NOTE: consider reorganizig the routing with more groups, this is getting kinda messy now

	user_router.GET("/profile", MiddlewareNoCache(), HandleGetProfile(&db))
	user_router.GET("/login", HandleGetLogin())
	user_router.POST("/login", HandlePostLogin(&db))
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup", HandleGetSignup())
	user_router.POST("/signup/send-mail", HandlePostSignup(&db))
	user_router.GET("/signup/from-mail/:token_id/:email", HandleGetSignupFromMail())
	user_router.POST("/signup/from-mail/:token_id/:email", HandlePostSignupFromMail(&db))
	user_router.GET("/delete_account", HandleGetDeleteAccount())
	user_router.POST("/delete_account", HandlePostDeleteAccount(&db))
	user_router.GET("/edit_profile", HandleGetEditProfile(&db))
	user_router.POST("/edit_profile", HandlePostEditProfile(&db))
	user_router.GET("/change_password", HandleGetChangePassword(&db))
	user_router.GET("/change_password/:token_id/:email", HandleGetChangePasswordFromMail())
	user_router.POST("/change_password/:tokenid/:email", HandlePostChangePasswordFromMail(&db))
	user_router.GET("/create_plan", HandleGetCreatePlan(&db))
	user_router.POST("/create_plan", HandlePostCreatePlan(&db))
	user_router.GET("/profile/plans/view/current", HandleGetViewCurrentPlan(&db))
	user_router.GET("/profile/plans/view/:wp_id", HandleGetViewPlan(&db))
	user_router.GET("/profile/plans/all_plans", HandleGetViewAllUserPlans(&db))
	user_router.GET("/profile/plans/make_current/:wp_id", HandleGetMakePlanCurrent(&db))
	user_router.GET("/forgot_password", HandleGetForgotPassword())
	user_router.POST("/forgot_password", HandlePostForgotPassword(&db))
	user_router.GET("/forgot_password/from-mail/:token_id/:email", HandleGetChangePassFromMail())
	user_router.GET("/tracks/view", HandleGetTracks(&db))
	user_router.GET("/tracks/view/:user_id", HandleGetTracks(&db)) // NOTE: this is for when someone is looking at another persons tracks
	user_router.GET("/tracks/view/:user_id/:track_id", HandleGetViewTrack(&db)) // NOTE: perhaps store the last updated track in a cookie to allow the user a quick access to the tracking page
	user_router.GET("/tracks/create", HandleGetTracksCreate(&db))
	user_router.POST("/tracks/create/:plan_id", HandlePostTracksCreate(&db))
	user_router.GET("/tracks/edit/:user_id/:track_id", HandleGetTracksEdit(&db))
	user_router.POST("/tracks/edit/:user_id/:track_id", HandlePostTracksEdit(&db))

	handler := sessionManager.LoadAndSave(router)
	handler = csrf.Protect(
		[]byte(csrf_key),
		csrf.Secure(true),
	)(handler)

	http.ListenAndServe(":8080", handler)
}
