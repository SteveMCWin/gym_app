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
		"add": func(a, b int) int {
			return a + b
		},
    }
}

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

	err = db.CacheAllExercises()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllTargets()
	if err != nil {
		panic(err)
	}

	err = db.LinkCachedExercisesAndTargets()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllPlansBasic()
	if err != nil {
		panic(err)
	}

	err = db.CacheAllGyms()
	if err != nil {
		panic(err)
	}

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
	router.GET("/ping", HandleGetPing())

	user_router := router.Group("/user")
	plan_router := router.Group("/user/:id/plan")
	track_router := router.Group("/user/:id/track")
	gym_router := router.Group("/gym")

	user_router.GET("/profile", HandleGetCurrentProfile())
	user_router.GET("/:id", MiddlewareNoCache(), HandleGetProfile(&db))
	user_router.GET("/search", HandleGetSearchForUser(&db))
	user_router.GET("/login", HandleGetLogin())
	user_router.POST("/login", HandlePostLogin(&db))
	user_router.GET("/logout", HandleGetLogout(&db))
	user_router.GET("/signup", HandleGetSignup())
	user_router.POST("/signup/send-mail", HandlePostSignup(&db))
	user_router.GET("/signup/from-mail/:token_id/:email", MiddlewareNoCache(), HandleGetSignupFromMail()) // Gotta check if the middleware no chache is actually doing it's job here
	user_router.POST("/signup/from-mail/:token_id/:email", HandlePostSignupFromMail(&db))
	user_router.GET("/delete_account", HandleGetDeleteAccount())
	user_router.POST("/delete_account", HandlePostDeleteAccount(&db))
	user_router.GET("/edit_profile", HandleGetEditProfile(&db))
	user_router.POST("/edit_profile", HandlePostEditProfile(&db))
	user_router.GET("/change_password", HandleGetChangePassword(&db))
	user_router.GET("/change_password/:token_id/:email", HandleGetChangePasswordFromMail())
	user_router.POST("/change_password/:tokenid/:email", HandlePostChangePasswordFromMail(&db))
	user_router.GET("/forgot_password", HandleGetForgotPassword())
	user_router.POST("/forgot_password", HandlePostForgotPassword(&db))
	user_router.GET("/forgot_password/from-mail/:token_id/:email", HandleGetChangePassFromMail())

	plan_router.GET("/create", HandleGetCreatePlan())
	plan_router.POST("/create", HandlePostCreatePlan(&db))
	plan_router.GET("/view/:wp_id", HandleGetViewPlan(&db))
	plan_router.GET("/view_all", HandleGetViewAllUserPlans(&db))
	plan_router.GET("/make_current/:wp_id", HandleGetMakePlanCurrent(&db))
	plan_router.GET("/edit/:wp_id", HandleGetEditPlan(&db))
	plan_router.POST("/edit/:wp_id", HandlePostEditPlan(&db))
	plan_router.GET("/get_plan_json/:wp_id", HandleGetPlanJSON(&db))

	track_router.GET("/view_all", HandleGetTracks(&db))
	track_router.GET("/view/latest", HandleGetTracksViewLatest(&db))
	track_router.GET("/view/:track_id", HandleGetViewTrack(&db))
	track_router.GET("/create", HandleGetTracksCreate(&db))
	track_router.POST("/create/:plan_id", HandlePostTracksCreate(&db))
	track_router.GET("/edit/:track_id", HandleGetTracksEdit(&db))
	track_router.POST("/edit/:track_id", HandlePostTracksEdit(&db))
	track_router.GET("/delete/:track_id", HandleGetTracksDelete(&db))

	gym_router.GET("/view_all", HandleGetViewAllGyms())
	gym_router.GET("/view/:gym_id", HandleGetViewGym(&db))

	handler := sessionManager.LoadAndSave(router)
	handler = csrf.Protect(
		[]byte(csrf_key),
		csrf.Secure(true),
	)(handler)

	http.ListenAndServe(":8080", handler)
}
