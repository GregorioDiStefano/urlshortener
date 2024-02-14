package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	db *sql.DB
}

type Database interface {
	GetUser(email string) (*User, error)
	SignupUser(email, password string) error

	InsertURL(userID int, url string) (int64, string, error)
	GetURL(id uint64) (string, string, error)
	GetURLs(userID int) ([]UserURLs, error)
	DisableURL(userID int, dbID uint64) error
	UpdateAccessAndLastAccessed(id uint64) error

	ValidateUser(email string) error

	Ping() error
	GetConnection() *sql.DB
	Close() error
}

type User struct {
	id    int
	email string

	password_hash []byte
}

type URL struct {
	Key          string     `json:"key" binding:"required"`
	Target       string     `json:"target" binding:"required"`
	Nonce        string     `json:"nonce" binding:"required"`
	Created      time.Time  `json:"created" binding:"required"`
	UserID       int        `json:"user_id" binding:"required"`
	LastAccessed *time.Time `json:"last_accessed" binding:"required"`
	AccessCount  int        `json:"access_count" binding:"required"`
}

// UserURLs are
type UserURLs struct {
	ShortURL     string  `json:"short_url"`
	Target       string  `json:"target"`
	Created      string  `json:"created"`
	LastAccessed *string `json:"last_accessed"`
	AccessCount  int     `json:"access_count"`
}

type dbConfig struct {
	host     string
	port     string
	user     string
	password string
	dbname   string
}

func NewDB(config *dbConfig) (Database, error) {
	// NOTE: https://go.dev/doc/database/sql-injection
	// Prepared statements are used, so don't worry about SQL injection
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.host, config.port, config.user, config.password, config.dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	// create tables if they don't exist, i would keep this here in a real world scenario, but
	// for the sake of the exercise, i'll leave it here
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL UNIQUE,
		password_hash BYTEA NOT NULL,
		created DATE NOT NULL,
		validate BOOLEAN NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS short_urls (
		id SERIAL PRIMARY KEY,
		target VARCHAR(255) NOT NULL,
		nonce VARCHAR(2) NOT NULL,
		created DATE NOT NULL,
		user_id INT,
		last_accessed DATE,
		access_count INT,
		disabled BOOLEAN DEFAULT FALSE,
		FOREIGN KEY (user_id) REFERENCES users(id)
	);

	CREATE TABLE IF NOT EXISTS schema_migrations (
		migration_id VARCHAR(255) PRIMARY KEY,
		migration_name VARCHAR(255) NOT NULL,
		applied_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	SELECT setval('users_id_seq', 999);
	`)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) GetURLs(userID int) ([]UserURLs, error) {
	var urls []UserURLs
	rows, err := d.db.Query(
		"SELECT id, nonce, target, created, last_accessed, access_count FROM short_urls WHERE user_id = $1 AND disabled = false", userID)

	fmt.Println(rows, err)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var url UserURLs

		var id int64
		var nonce string

		err := rows.Scan(&id, &nonce, &url.Target, &url.Created, &url.LastAccessed, &url.AccessCount)
		if err != nil {
			return nil, err
		}

		// construct the short url
		url.ShortURL = fmt.Sprintf("%s%s", uint64ToBase64(int64(id)), nonce)
		urls = append(urls, url)
	}

	return urls, nil
}

// GetUser returns a user by email
func (d *DB) GetUser(email string) (*User, error) {
	var user User

	err := d.db.QueryRow("SELECT id, password_hash, email FROM users WHERE email = $1", email).Scan(
		&user.id, &user.password_hash, &user.email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DB) SignupUser(email, password string) error {
	var existingUsername string
	err := d.db.QueryRow("SELECT email FROM users WHERE email = $1", email).Scan(&existingUsername)

	if err == sql.ErrNoRows {
		// Username does not exist, so it's safe to proceed with creating a new user
		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}

		timeNow := time.Now().Format("2006-01-02")
		result, err := d.db.Exec(
			"INSERT INTO users (email, password_hash, created, validate) VALUES ($1, $2, $3, $4)",
			email,
			bcryptPassword,
			timeNow,
			true, // setting to true, would be false if we did the validation procedure
		)

		if err != nil {
			return fmt.Errorf("error creating user: %w", err)
		}

		if changed, err := result.RowsAffected(); err != nil || changed != 1 {
			return fmt.Errorf("error creating user: %w", err)
		}

		return err
	} else if err != nil {
		return fmt.Errorf("error checking email existence: %w", err)
	} else {
		return fmt.Errorf("email already exists")
	}
}

func (d *DB) InsertURL(userID int, url string) (int64, string, error) {
	nonce := randomString(2)

	now := time.Now().Format("2006-01-02")
	var id int64

	err := d.db.QueryRow("INSERT INTO short_urls (target, created, user_id, last_accessed, access_count, nonce) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		url, now, userID, nil, 0, nonce).Scan(&id)

	return id, nonce, err
}

func (d *DB) UpdateAccessAndLastAccessed(id uint64) error {
	_, err := d.db.Exec("UPDATE short_urls SET access_count = access_count + 1, last_accessed = NOW() WHERE id = $1", id)
	return err
}

func (d *DB) GetURL(id uint64) (string, string, error) {
	// increment value and return value from column
	var target string
	var nonce string

	err := d.db.QueryRow(
		"SELECT target, nonce FROM short_urls WHERE id = $1 AND disabled = false", id).Scan(&target, &nonce)
	return target, nonce, err
}

func (d *DB) DisableURL(userID int, dbID uint64) error {
	// update user to be validated
	result, err := d.db.Exec("UPDATE short_urls SET disabled = true WHERE user_id = $1 AND id = $2", userID, dbID)

	if err != nil {
		return fmt.Errorf("error disabling url: %w", err)
	}

	if changed, err := result.RowsAffected(); err != nil || changed != 1 {
		return fmt.Errorf("error disabling url: %w", err)
	} else {
		return err
	}
}

func (d *DB) ValidateUser(email string) error {
	// update user to be validated
	result, err := d.db.Exec("UPDATE users SET validate = true WHERE email = $1", email)

	if err != nil {
		return fmt.Errorf("error validating user: %w", err)
	}

	if changed, err := result.RowsAffected(); err != nil || changed != 1 {
		return fmt.Errorf("error validating user: %w", err)
	} else {
		return err
	}
}

func (d *DB) Ping() error {
	return d.db.Ping()
}

func (d *DB) GetConnection() *sql.DB {
	return d.db
}

func randomString(size int) string {
	b := make([]byte, size+2)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[2 : size+2]
}
