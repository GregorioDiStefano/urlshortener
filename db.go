package main

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"time"
	"unsafe"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type DB struct {
	db *sql.DB
}

type Database interface {
	GetUser(username string) (*User, error)
	SignupUser(username, password, email string) error

	InsertURL(userID int, url string) (int64, string, error)
	GetURL(id uint64) (string, error)
	GetURLs(userID int) ([]URL, error)

	ValidateUser(username string) error

	Ping() error
}

type User struct {
	id       int
	username string
	email    string

	password_hash []byte
}

type URL struct {
	id           int
	Key          string     `json:"key" binding:"required"`
	Target       string     `json:"target" binding:"required"`
	Nonce        string     `json:"nonce" binding:"required"`
	Created      time.Time  `json:"created" binding:"required"`
	UserID       int        `json:"user_id" binding:"required"`
	LastAccessed *time.Time `json:"last_accessed" binding:"required"`
	AccessCount  int        `json:"access_count" binding:"required"`
}

func NewDB() (Database, error) {
	// NOTE: https://go.dev/doc/database/sql-injection
	// Prepared statements are used, so don't worry about SQL injection
	db, err := sql.Open("postgres", "user=postgres password=mysecretpassword sslmode=disable")

	if err != nil {
		return nil, err
	}

	// create tables if they don't exist, i would keep this here in a real world scenario, but
	// for the sake of the exercise, i'll leave it here
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		username VARCHAR(255) NOT NULL UNIQUE,
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
		FOREIGN KEY (user_id) REFERENCES users(id)
	);
	
	ALTER SEQUENCE users_id_seq RESTART WITH 1000;
	ALTER SEQUENCE short_urls_id_seq RESTART WITH 1000;
	
	`)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func (d *DB) GetURLs(userID int) ([]URL, error) {
	var urls []URL
	rows, err := d.db.Query("SELECT id, target, nonce, created, user_id, last_accessed, access_count FROM short_urls WHERE user_id = $1", userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var url URL
		err := rows.Scan(&url.id, &url.Target, &url.Nonce, &url.Created, &url.UserID, &url.LastAccessed, &url.AccessCount)

		fmt.Println(url)
		if err != nil {
			return nil, err
		}

		url.Key = idToKey(int64(url.id))
		urls = append(urls, url)

	}

	return urls, nil
}

// GetUser returns a user by username
func (d *DB) GetUser(username string) (*User, error) {
	var user User

	err := d.db.QueryRow("SELECT id, username, password_hash, email FROM users WHERE username = $1", username).Scan(
		&user.id, &user.username, &user.password_hash, &user.email)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (d *DB) SignupUser(username, password, email string) error {
	var existingUsername string
	err := d.db.QueryRow("SELECT username FROM users WHERE username = $1", username).Scan(&existingUsername)

	if err == sql.ErrNoRows {
		// Username does not exist, so it's safe to proceed with creating a new user
		bcryptPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

		if err != nil {
			return fmt.Errorf("error hashing password: %v", err)
		}

		timeNow := time.Now().Format("2006-01-02")
		result, err := d.db.Exec(
			"INSERT INTO users (username, password_hash, email, created, validate) VALUES ($1, $2, $3, $4, $5)",
			username,
			bcryptPassword,
			email,
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
		return fmt.Errorf("error checking username existence: %w", err)
	} else {
		return fmt.Errorf("username already exists")
	}
}

func (d *DB) InsertURL(userID int, url string) (int64, string, error) {
	nonce := randomString(2)

	now := time.Now().Format("2006-01-02")
	fmt.Println(now)
	var id int64

	err := d.db.QueryRow("INSERT INTO short_urls (target, created, user_id, last_accessed, access_count, nonce) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id",
		url, now, userID, nil, 0, nonce).Scan(&id)

	return id, nonce, err
}

func (d *DB) GetURL(id uint64) (string, error) {
	// increment value and return value from column
	var target string
	err := d.db.QueryRow(
		"UPDATE short_urls SET access_count = access_count + 1, last_accessed = NOW() WHERE id = $1 RETURNING target", id).Scan(&target)
	return target, err
}

func (d *DB) ValidateUser(username string) error {
	// update user to be validated
	result, err := d.db.Exec("UPDATE users SET validate = true WHERE username = $1", username)

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

func randomString(size int) string {
	var alphabet = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]byte, size)
	rand.Read(b)
	for i := 0; i < size; i++ {
		b[i] = alphabet[b[i]%byte(len(alphabet))]
	}
	return *(*string)(unsafe.Pointer(&b))
}
