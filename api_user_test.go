package main

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestInvalidShortURLs(t *testing.T) {
	app := NewTestApp()
	defer wipeDB(app.db)

	router := setupRouter(app)
	token := createUserAndGetToken(app)

	// Test invalid URL
	requestString := `{"url":"not a url"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/api/v1/shorten", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	assert.Equal(t, 400, loginW.Code)

	// Test missing URL
	requestString = `{}`
	buf = bytes.NewBufferString(requestString)
	req, _ = http.NewRequest("POST", "/api/v1/shorten", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	loginW = httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	assert.Equal(t, 400, loginW.Code)
}

// TestUserShortURLs tests that a user can create and retrieve short URLs
func TestUserShortURLs(t *testing.T) {
	app := NewTestApp()
	defer wipeDB(app.db)

	router := setupRouter(app)
	token := createUserAndGetToken(app)

	urls := []string{"https://www.google.com", "https://www.yahoo.com", "https://www.bing.com"}
	shortURLKeys := map[string]string{}

	// create short urls for the logged in user
	for _, url := range urls {
		requestString := `{"url":"` + url + `"}`
		buf := bytes.NewBufferString(requestString)
		req, _ := http.NewRequest("POST", "/api/v1/shorten", buf)
		req.Header.Set("Authorization", "Bearer "+token)
		loginW := httptest.NewRecorder()
		router.ServeHTTP(loginW, req)

		key := gjson.Get(loginW.Body.String(), "key").String()
		shortURLKeys[url] = key
	}

	// get the short urls for the logged in user
	req, _ := http.NewRequest("GET", "/api/v1/urls/", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	// check the output is what you woud expect
	assert.Equal(t, 200, loginW.Code)

	assert.Equal(t, "https://www.google.com", gjson.Get(loginW.Body.String(), "0.target").String())
	assert.Equal(t, "https://www.yahoo.com", gjson.Get(loginW.Body.String(), "1.target").String())
	assert.Equal(t, "https://www.bing.com", gjson.Get(loginW.Body.String(), "2.target").String())

	// Verify each short URL redirects correctly
	for target, key := range shortURLKeys {
		req, _ := http.NewRequest("GET", "/"+key, nil)
		rec := httptest.NewRecorder()
		router.ServeHTTP(rec, req)

		// Check the Location header
		location := rec.Header().Get("Location")
		assert.Equal(t, target, location)
	}
}

// TODO: Test accessing a short URL increases access count and time
// Test disable URL
func TestDisableURL(t *testing.T) {
	app := NewTestApp()
	defer wipeDB(app.db)

	router := setupRouter(app)
	token := createUserAndGetToken(app)

	// create short url
	requestString := `{"url":"https://www.google.com"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/api/v1/shorten", buf)
	req.Header.Set("Authorization", "Bearer "+token)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	fmt.Println("got: ", loginW.Body.String())
	key := gjson.Get(loginW.Body.String(), "key").String()

	// disable short url
	req, _ = http.NewRequest("DELETE", "/api/v1/shorten/"+key, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	loginW = httptest.NewRecorder()
	router.ServeHTTP(loginW, req)

	// check the output is what you woud expect
	assert.Equal(t, 200, loginW.Code)

	// Verify the short URL is disabled
	req, _ = http.NewRequest("GET", "/"+key, nil)
	rec := httptest.NewRecorder()
	router.ServeHTTP(rec, req)

	// Check the Location header
	location := rec.Header().Get("Location")
	assert.Equal(t, "", location)

	// and 404 is returned
	assert.Equal(t, 404, rec.Code)
}

// Test that a user can't see another user's short URLs
// TODO: implement this test
