package handlers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"time"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	// "github.com/stretchr/testify/require"

	"fitness_app/handlers"
	"fitness_app/models"
)

var db models.DataBase

var domain string
var csrf_key string
var handler http.Handler

func init() {

	_, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory:", err)
	}

	err = os.Chdir("..")
	if err != nil {
		log.Fatal("Couldn't change to project root:", err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		log.Fatal("Couldn't open DataBase, error:", err)
	}

	domain = os.Getenv("DOMAIN")
	csrf_key = os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		log.Fatal("Couldn't load .env variables")
	}

	gin.SetMode(gin.TestMode)

	err = db.InitDatabase(true)
	if err != nil {
		log.Fatal("Couldn't open DataBase, error:", err)
	}


	handler = handlers.SetUpRouter(domain, csrf_key, db)
}

func TestHandleGetProfile_NotLoggedIn(t *testing.T) {
	test_user_id := 4

	req := httptest.NewRequest(http.MethodGet, "/user/"+strconv.Itoa(test_user_id), nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusPermanentRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")
}

func TestHandleGetLogIn(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/user/login", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestHandlePostLogIn(t *testing.T) {
	g, ok := models.FetchCachedGym(3)
	if !ok {
		t.Error("Couldnt load gym with id 3")
	}

	test_user := models.User {
		Name: "TestUser",
		Email: "test@gmail.com",
		Password: "right_password",
		TrainingSince: time.Now(),
		IsTrainer: false,
		GymGoals: "strength",
		CurrentGym: *g,
		DateCreated: time.Now(),
	}

	user_id, err := db.CreateUser(test_user)
	if err != nil {
		t.Error("Couldn't create user!")
	}

	defer func() {
		_, err = db.DeleteUser(user_id)
		if err != nil {
			t.Error("Couldn't delete user!")
		}
	}()

	form := url.Values{}
	form.Add("email", "test@gmail.com")
	form.Add("password", "wrong_password")

	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")

	log.Println("Should have said 'wrong password' just now")

	form = url.Values{}
	form.Add("email", "test@gmail.com")
	form.Add("password", "right_password")

	req = httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(user_id))

}
