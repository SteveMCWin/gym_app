package handlers_test

import (
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

func TestHandleGetProfile_NotLoggedIn(t *testing.T) {

	dir, err := os.Getwd()
	if err != nil {
		t.Fatal("Couldn't get current working directory:", err)
	}

	err = os.Chdir("..")
	if err != nil {
		t.Fatal("Couldn't change to project root:", err)
	}

	defer func() {
		err = os.Chdir(dir)
		if err != nil {
			t.Fatal("couldn't change to starting dir")
		}
	}()

	err = godotenv.Load(".env")
	if err != nil {
		t.Error("Couldn't open DataBase, error:", err)
	}

	domain := os.Getenv("DOMAIN")
	csrf_key := os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		t.Error("Couldn't load .env variables")
	}

	gin.SetMode(gin.TestMode)
	//
	var db models.DataBase
	err = db.InitDatabase(true)
	if err != nil {
		t.Error("Couldn't open DataBase, error:", err)
	}

	handler := handlers.SetUpRouter(domain, csrf_key, db)
	req := httptest.NewRequest(http.MethodGet, "/user/4", nil)
	resp := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(resp)
	c.Request = req

	handler.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusPermanentRedirect, resp.Code)
	assert.Contains(t, resp.Header().Get("Location"), "/user/login")
}

func TestHandleGetProfile_LoggedIn(t *testing.T) {
	err := os.Chdir("..")
	if err != nil {
		t.Fatal("Couldn't change to project root:", err)
	}

	err = godotenv.Load(".env")
	if err != nil {
		t.Error("Couldn't load .env file:", err)
	}

	domain := os.Getenv("DOMAIN")
	csrf_key := os.Getenv("CSRF_KEY")

	if domain == "" || csrf_key == "" {
		t.Error("Couldn't load .env variables")
	}

	gin.SetMode(gin.TestMode)

	var db models.DataBase
	err = db.InitDatabase(true)
	if err != nil {
		t.Error("Couldn't open DataBase, error:", err)
	}

	log.Println("This is fine 1")

	test_user_id := 4
	
	handler := handlers.SetUpRouter(domain, csrf_key, db)
	req := httptest.NewRequest(http.MethodGet, "/user/"+strconv.Itoa(test_user_id), nil)
	resp := httptest.NewRecorder()

	log.Println("This is fine 2")
	c, _ := gin.CreateTestContext(resp)

	handlers.SessionManager.Put(c.Request.Context(), "user_id", test_user_id)

	c.Request = req

	handler.ServeHTTP(resp, req)

	log.Println("This is fine 4")
	assert.Equal(t, http.StatusOK, resp.Code)
}
