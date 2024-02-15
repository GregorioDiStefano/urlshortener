package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"

	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

func newDBConnection() (Database, error) {
	dbConfig := &dbConfig{
		host:     "localhost",
		port:     "9999",
		user:     "testuser",
		password: "testpassword",
		dbname:   "testdb",
	}

	return NewDB(dbConfig)
}

func wipeDB(db Database) {
	_, err := db.GetConnection().Exec("DROP TABLE IF EXISTS short_urls; DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS schema_migrations")
	if err != nil {
		log.Fatal(err)
	}
}

func NewTestApp() *App {
	db, err := newDBConnection()

	if err != nil {
		log.Fatal(err)
	}

	cache, err := NewCache("localhost")
	if err != nil {
		log.Fatal(err)
	}

	tokenChecker := jwtAuthMiddleware{secret: []byte("testsecret1234")}

	return &App{
		db:              db,
		cache:           cache,
		tokenMiddleware: tokenChecker,
	}
}

// crateUserAndGetToken creates a user and returns the token, used for testing
func createUserAndGetToken(app *App) string {
	w := httptest.NewRecorder()

	requestString := `{"password":"testpassword!","email":"greg@greg.ca"}`
	buf := bytes.NewBufferString(requestString)
	req, _ := http.NewRequest("POST", "/auth/register", buf)
	router := setupRouter(app)
	router.ServeHTTP(w, req)

	loginRequestString := `{"email":"greg@greg.ca","password":"testpassword!"}`
	loginBuf := bytes.NewBufferString(loginRequestString)
	loginReq, _ := http.NewRequest("POST", "/auth/login", loginBuf)
	loginW := httptest.NewRecorder()
	router.ServeHTTP(loginW, loginReq)

	return gjson.Get(loginW.Body.String(), "token").String()
}
