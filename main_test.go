package main

import (
	"fmt"
	"testing"

	log "github.com/sirupsen/logrus"
)

// Drop tables before tests and create new connection
func TestMain(m *testing.M) {
	log.SetLevel(log.DebugLevel)
	db, err := newDBConnection()

	if err != nil {
		log.Fatal(err)
	}

	_, err = db.GetConnection().Exec("DROP TABLE IF EXISTS short_urls; DROP TABLE IF EXISTS users; DROP TABLE IF EXISTS schema_migrations")
	fmt.Println("Dropped tables")
	if err != nil {
		log.Fatal(err)
	}

	m.Run()
}
