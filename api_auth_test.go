package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterOK(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email": "greg@greg.com", "password":"testpassword!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
}

func TestRegisterDuplicate(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email": "greg@greg.com", "password":"testpassword!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// attempt to register again
	w = httptest.NewRecorder()
	buf = bytes.NewBufferString(requestString)
	req, _ = http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 500, w.Code)
	assert.Equal(t, `{"error":"email already exists"}`, w.Body.String())
}

func TestRegisterBadJSON(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Equal(t, `{"error":"invalid request: unexpected EOF"}`, w.Body.String())
}

func TestRegisterBadPassword(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email": "greg@greg.com", "password":"test!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	assert.Contains(t, w.Body.String(), "failed on the 'min' tag")
}

func TestLoginBadAccount(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email":"greg@foobar.com","password":"testpassword!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/login", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestLoginBadPassword(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email":"greg@greg.com","password":"testpassword!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// attempt to login
	loginRequestString := `{"email":"greg@greg.com", "password":"badpassword"}`
	loginBuf := bytes.NewBufferString(loginRequestString)
	loginReq, _ := http.NewRequest("POST", "/auth/login", loginBuf)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	assert.Equal(t, 401, loginW.Code)
	assert.Equal(t, `{"error":"invalid credentials"}`, loginW.Body.String())
}

func TestSuccessfulRegisterWithLogin(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	w := httptest.NewRecorder()

	requestString := `{"email": "greg@greg.com", "password":"testpassword!"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)

	// attempt to login
	loginRequestString := `{"email":"greg@greg.com","password":"testpassword!"}`
	loginBuf := bytes.NewBufferString(loginRequestString)
	loginReq, _ := http.NewRequest("POST", "/auth/login", loginBuf)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)
	assert.Equal(t, 200, loginW.Code)
}

func TestInvalidToken(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	badTokens := []string{"", "Bearer ", "Bearer invalid-token"}
	requestString := `{"url":"https://www.google.com"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/api/v1/shorten", buf)

	for _, token := range badTokens {
		req.Header.Set("Authorization", token)
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, req)

		assert.Equal(t, 401, loginW.Code)
	}
}

func TestMissingToken(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)
	defer wipeDB(app.db)

	requestString := `{"url":"https://www.google.com"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/api/v1/shorten", buf)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	assert.Equal(t, 401, loginW.Code)
	assert.Equal(t, `{"error":"authorization header is required"}`, loginW.Body.String())
}
