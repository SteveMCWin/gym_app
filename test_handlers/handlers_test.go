package handlers_test

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"fitness_app/handlers"
	"fitness_app/models"
)

var db models.DataBase

var domain string
var csrf_key string
var handler http.Handler

const (
	TEST_USER_EMAIL = "test@gmail.com"
	TEST_USER_PASS  = "right_password"
)

var test_user_id int

func init() {

	// NOTE: the env vars can be read without loading the env when running the test in github actions, so first check that
	if os.Getenv("CI") == "" {
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
			log.Fatal("Couldn't load the .env:", err)
		}
	}

	domain = os.Getenv("DOMAIN")
	csrf_key = os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		log.Fatal("Couldn't load .env variables")
	}

	gin.SetMode(gin.TestMode)

	err := db.InitDatabase(true)
	if err != nil {
		log.Fatal("Couldn't open DataBase, error:", err)
	}

	err = db.CacheData()
	if err != nil {
		log.Fatal("Couldn't cache data, error:", err)
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

func TestCreateUser(t *testing.T) {

	// If user exists, don't create him again
	user_id, err := db.AuthUserByEmail(TEST_USER_EMAIL, TEST_USER_PASS)
	if err == nil {
		test_user_id = user_id
		log.Println("there is a user with email:", TEST_USER_EMAIL)
		return
	}

	g, ok := models.FetchCachedGym(3)
	if !ok {
		t.Error("Couldnt load gym with id 3")
	}

	test_user := models.User{
		Name:          "TestUser",
		Email:         TEST_USER_EMAIL,
		Password:      TEST_USER_PASS,
		TrainingSince: time.Now(),
		IsTrainer:     false,
		GymGoals:      "strength",
		CurrentGym:    *g,
		DateCreated:   time.Now(),
	}

	// Creating test user
	user_id, err = db.CreateUser(test_user)
	if err != nil {
		t.Error("Couldn't create user!")
	}

	test_user_id = user_id

	_, err = db.ReadUser(user_id)
	if err != nil {
		t.Error("Couldn't read the user data")
	}

	// defer func() {
	// 	_, err = db.DeleteUser(user_id)
	// 	if err != nil {
	// 		t.Error("Couldn't delete user!")
	// 	}
	// }()

}

func TestHandlePostLogIn(t *testing.T) {

	// Testing login with wrong password
	form := url.Values{}
	form.Add("email", TEST_USER_EMAIL)
	form.Add("password", "wrong_password")

	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")

	// Testing login with correct password

	form = url.Values{}
	form.Add("email", TEST_USER_EMAIL)
	form.Add("password", TEST_USER_PASS)

	req = httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	// When the user is logged in succesfully, he should get a session cookie that keeps him logged in
	var sessionCookie *http.Cookie
	for _, cookie := range resp.Result().Cookies() {
		if strings.Contains(cookie.Name, "session") {
			sessionCookie = cookie
			break
		}
	}

	assert.NotNil(t, sessionCookie, "Session cookie should be set after login")

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(test_user_id))

	// Once the user is logged in, going to the login page should redirect him to the profile page
	req = httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(test_user_id))

	// Once the user is logged in, going to the signup page should redirect him to the profile page
	req = httptest.NewRequest(http.MethodGet, "/user/signup", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(test_user_id))

	// Logging out should remove the cookie and redirect the user to the login page
	req = httptest.NewRequest(http.MethodGet, "/user/logout", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")

	// The cookie should not work anymore since it's deleted in the session manager
	req = httptest.NewRequest(http.MethodGet, "/user/login", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestCreatePlan(t *testing.T) {

	plan_json := `{
	"name":"Ime plana",
	"description":"ovo je opciono valjda",
	"make_current":true,
	"days":[{"name":"prvi","exercises":[{"exercise":{"name":"barbellbench press"},"weight":55,"unit":"kg","sets":3,"min_reps":5,"max_reps":12},{"exercise":{"name":"dumbbell curls"},"weight":15,"unit":"kg","sets":4,"min_reps":7,"max_reps":10}]},{"name":"drugi","exercises":[{"exercise":{"name":"t-bar row"},"weight":40,"unit":"kg","sets":3,"min_reps":7,"max_reps":10}]},{"name":"treci","exercises":[{"exercise":{"name":"hack sqaut"},"weight":3,"unit":"kg","sets":65,"min_reps":6,"max_reps":12},{"exercise":{"name":"rear delt row"},"weight":15,"unit":"kg","sets":5,"min_reps":10,"max_reps":15}]}]
	}
	`

	// Log the user in
	form := url.Values{}
	form.Add("email", TEST_USER_EMAIL)
	form.Add("password", TEST_USER_PASS)

	req := httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	// When the user is logged in succesfully, he should get a session cookie that keeps him logged in
	var sessionCookie *http.Cookie
	for _, cookie := range resp.Result().Cookies() {
		if strings.Contains(cookie.Name, "session") {
			sessionCookie = cookie
			break
		}
	}

	assert.NotNil(t, sessionCookie, "Session cookie should be set after login")

	req = httptest.NewRequest(http.MethodPost, "/user/"+strconv.Itoa(test_user_id)+"/plan/create", strings.NewReader(plan_json))
	req.Header.Set("Content-Type", "application/json")
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
	}

	log.Println("Response body:", string(body))
	wp_id_str := string(body)

	req = httptest.NewRequest(http.MethodGet, "/user/"+strconv.Itoa(test_user_id)+"/plan/view/"+wp_id_str, nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

}

func TestDeleteAccount(t *testing.T) {

	// The user shouldn't be able to delete their account if they are not logged in
	req := httptest.NewRequest(http.MethodGet, "/user/delete_account", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusTemporaryRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")

	// Log in user
	form := url.Values{}
	form.Add("email", TEST_USER_EMAIL)
	form.Add("password", TEST_USER_PASS)

	req = httptest.NewRequest(http.MethodPost, "/user/login", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	// When the user is logged in succesfully, he should get a session cookie that keeps him logged in
	var sessionCookie *http.Cookie
	for _, cookie := range resp.Result().Cookies() {
		if strings.Contains(cookie.Name, "session") {
			sessionCookie = cookie
			break
		}
	}

	assert.NotNil(t, sessionCookie, "Session cookie should be set after login")

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(test_user_id))

	// Logged out user and correct password
	req = httptest.NewRequest(http.MethodGet, "/user/delete_account", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	assert.Equal(t, http.StatusOK, resp.Code)

	form = url.Values{}
	form.Add("password", TEST_USER_PASS)

	req = httptest.NewRequest(http.MethodPost, "/user/delete_account", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// No session cookie
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")

	// Logged in user and wrong password
	req = httptest.NewRequest(http.MethodGet, "/user/delete_account", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	assert.Equal(t, http.StatusOK, resp.Code)

	form = url.Values{}
	form.Add("password", "not_the_right_password")

	req = httptest.NewRequest(http.MethodPost, "/user/delete_account", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/"+strconv.Itoa(test_user_id))

	// Logged in user and correct password
	req = httptest.NewRequest(http.MethodGet, "/user/delete_account", nil)
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	assert.Equal(t, http.StatusOK, resp.Code)

	form = url.Values{}
	form.Add("password", TEST_USER_PASS)

	req = httptest.NewRequest(http.MethodPost, "/user/delete_account", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(sessionCookie)
	resp = httptest.NewRecorder()

	c, _ = gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusSeeOther, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/")
}
