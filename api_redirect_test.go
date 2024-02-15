package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	mock "urlshort/mocks"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestCreateShortlinkOK(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	// Register a user
	requestString := `{"password":"testpassword!","email":"greg@greg.com"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// Login
	loginRequestString := `{"email":"greg@greg.com","password":"testpassword!"}`
	loginBuf := bytes.NewBufferString(loginRequestString)
	loginReq, _ := http.NewRequest("POST", "/auth/login", loginBuf)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)
	assert.Equal(t, 200, loginW.Code)

	token := gjson.Get(loginW.Body.String(), "token")

	requestString = `{"url":"https://www.google.com"}`
	buf = bytes.NewBufferString(requestString)
	req, _ = http.NewRequest("POST", "/api/v1/shorten", buf)
	req.Header.Set("Authorization", "Bearer "+token.String())
	loginW = httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	fmt.Println(loginW.Body.String())
	assert.Equal(t, 201, loginW.Code)

	key := gjson.Get(loginW.Body.String(), "key")

	req, _ = http.NewRequest("GET", "/"+key.String(), nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, 301, w.Code)
	assert.Equal(t, "https://www.google.com", w.Header().Get("Location"))
}

func TestRedisKeySet(t *testing.T) {
	app := NewTestApp()
	cacheMock := mock.NewCache(t)

	app.cache = cacheMock

	cacheMock.On("GetURL", "Aw8c").Return("https://google.ca", nil)

	// call /testkey and check that the cache is called
	router := setupRouter(app)
	req, _ := http.NewRequest("GET", "/Aw8c", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// check that the db is not called
}
