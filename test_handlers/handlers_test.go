package handlers_test

import (
	"fmt"
	// "errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	// "strings"
	"testing"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/joho/godotenv"
	// "github.com/stretchr/testify/require"

	"fitness_app/handlers"
	"fitness_app/models"
)

var db models.DataBase

var domain string
var csrf_key string

func init() {

	_, err := os.Getwd()
	if err != nil {
		log.Fatal("Couldn't get current working directory:", err)
	}

	err = os.Chdir("..")
	if err != nil {
		log.Fatal("Couldn't change to project root:", err)
	}

	// defer func() {
	// 	err = os.Chdir(dir)
	// 	if err != nil {
	// 		log.Fatal("couldn't change to starting dir")
	// 	}
	// }()

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
}

func TestHandleGetProfile_NotLoggedIn(t *testing.T) {
	log.Println("LoggedOut test")

	test_user_id := 4

	handler := handlers.SetUpRouter(domain, csrf_key, db)
	req := httptest.NewRequest(http.MethodGet, "/user/"+strconv.Itoa(test_user_id), nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusPermanentRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")
}

func TestHandleGetProfile_LoggedIn(t *testing.T) {
	log.Println("LoggedIn test")

	test_user_id := 4
	
	handler := handlers.SetUpRouter(domain, csrf_key, db)

	sessionSetupReq := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/user/%d", test_user_id), nil)
	sessionSetupResp := httptest.NewRecorder()
	
	handler.ServeHTTP(sessionSetupResp, sessionSetupReq)

	t.Logf("Session setup response code: %d", sessionSetupResp.Code)
	t.Logf("Session setup response body: %s", sessionSetupResp.Body.String())
	t.Logf("Session setup response headers: %v", sessionSetupResp.Header())

	var sessionCookie *http.Cookie
	for _, cookie := range sessionSetupResp.Result().Cookies() {
		t.Logf("fount cookie: %s", cookie.Name)
		if cookie.Name == "test_session" { // Check your SCS config for actual cookie name
			sessionCookie = cookie
			break
		}
	}
	
	if sessionCookie == nil {
		t.Fatal("No session cookie found after setting session")
	}
	
	// Step 2: Make the actual test request with the session cookie
	req := httptest.NewRequest(http.MethodGet, "/user/"+strconv.Itoa(test_user_id), nil)
	req.AddCookie(sessionCookie)
	
	resp := httptest.NewRecorder()
	handler.ServeHTTP(resp, req)
	
	// Test the response
	assert.Equal(t, http.StatusOK, resp.Code)
}
