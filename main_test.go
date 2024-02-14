package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

func TestPingRoute(t *testing.T) {
	app := NewTestApp()
	router := setupRouter(app)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, `{"message":"pong"}`, w.Body.String())
}

func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	db, err := newDBConnection()

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.GetConnection().Exec("DROP TABLE IF EXISTS short_urls; DROP TABLE IF EXISTS users;")
	fmt.Println("Dropped tables")
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
