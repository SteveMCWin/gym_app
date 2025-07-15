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

	// err = db.CacheAllUserNamesAndIds()
	// if err != nil {
	// 	panic(err)
	// }

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
	// plan_router := router.Group("/user/:id/plan")
	// track_router := router.Group("/user/track")

	user_router.GET("/:id", MiddlewareNoCache(), HandleGetProfile(&db))
	user_router.GET("/search", HandleGetSearchForUser(&db))
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
	user_router.GET("/forgot_password", HandleGetForgotPassword())
	user_router.POST("/forgot_password", HandlePostForgotPassword(&db))
	user_router.GET("/forgot_password/from-mail/:token_id/:email", HandleGetChangePassFromMail())

	user_router.GET("/create_plan", HandleGetCreatePlan())
	user_router.POST("/create_plan", HandlePostCreatePlan(&db))
	user_router.GET("/profile/plans/view/current", HandleGetViewCurrentPlan(&db))
	user_router.GET("/profile/plans/view/:wp_id", HandleGetViewPlan(&db))
	user_router.GET("/profile/plans/all_plans", HandleGetViewAllUserPlans(&db))
	user_router.GET("/profile/plans/make_current/:wp_id", HandleGetMakePlanCurrent(&db))
	user_router.GET("/profile/plans/edit/:wp_id", HandleGetEditPlan(&db))
	user_router.POST("/profile/plans/edit/:wp_id", HandlePostEditPlan(&db))
	user_router.GET("/get_plan_json/:wp_id", HandleGetPlanJSON(&db))

	// plan_router.GET("/create", HandleGetCreatePlan())
	// plan_router.POST("/create", HandlePostCreatePlan(&db))
	// plan_router.GET("/view/current", HandleGetViewCurrentPlan(&db))
	// plan_router.GET("/view/:wp_id", HandleGetViewPlan(&db))
	// plan_router.GET("/view_all", HandleGetViewAllUserPlans(&db))
	// plan_router.GET("/make_current/:wp_id", HandleGetMakePlanCurrent(&db))
	// plan_router.GET("/edit/:wp_id", HandleGetEditPlan(&db))
	// plan_router.POST("/edit/:wp_id", HandlePostEditPlan(&db))
	// plan_router.GET("/get_plan_json/:wp_id", HandleGetPlanJSON(&db))

	user_router.GET("/tracks/view", HandleGetTracks(&db))
	user_router.GET("/tracks/view/:user_id", HandleGetTracks(&db)) // NOTE: this is for when someone is looking at another persons tracks
	user_router.GET("/tracks/create", HandleGetTracksCreate(&db))
	user_router.POST("/tracks/create/:plan_id", HandlePostTracksCreate(&db))
	user_router.GET("/tracks/edit/:user_id/:track_id", HandleGetTracksEdit(&db))
	user_router.POST("/tracks/edit/:user_id/:track_id", HandlePostTracksEdit(&db))

	// track_router.GET("/view_all", HandleGetTracks(&db))
	// track_router.GET("/view/:track_id", HandleGetViewTrack(&db))
	// track_router.GET("/create", HandleGetTracksCreate(&db))
	// track_router.POST("/create/:plan_id", HandlePostTracksCreate(&db))
	// track_router.GET("/edit/:track_id", HandleGetTracksEdit(&db))
	// track_router.POST("/edit/:track_id", HandlePostTracksEdit(&db))

	handler := sessionManager.LoadAndSave(router)
	handler = csrf.Protect(
		[]byte(csrf_key),
		csrf.Secure(true),
	)(handler)

	http.ListenAndServe(":8080", handler)
}
