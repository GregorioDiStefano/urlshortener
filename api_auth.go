package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (app *App) login(c *gin.Context) {
	type Login struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		log.WithError(err).Error("error binding json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *User
	var err error
	if user, err = app.db.GetUser(json.Email); err != nil {
		log.WithError(err).Error("error getting user")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword(user.passwordHash, []byte(json.Password)); err != nil {
		log.WithError(err).Error("error comparing password")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := app.tokenMiddleware.GenerateUserToken(user.id)
	if err != nil {
		log.WithError(err).Error("error generating token")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (app *App) register(c *gin.Context) {
	type RegisterRequest struct {
		Password string `json:"password" binding:"required,ascii,min=8,max=255"`
		Email    string `json:"email" binding:"required,email"`
	}

	validatePassword := func(login *RegisterRequest) error {
		const specialChars = "!@#$%^&*()_+"
		if !strings.ContainsAny(login.Password, specialChars) {
			return fmt.Errorf("password must contain at least one special character")
		}

		return nil
	}

	var json RegisterRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		log.WithError(err).Error("error binding json")
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Errorf("invalid request: %w", err).Error()})
		return
	}

	if err := validatePassword(&json); err != nil {
		log.WithError(err).Error("error validating")
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed: " + err.Error()})
		return
	}

	if err := app.db.SignupUser(json.Email, json.Password); err != nil {
		log.WithError(err).Error("error signing up user")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Note: usually you would want to valid the user's email address before allowing them to login.
	// You can do this by signing a JWT token with the user's email address and sending them a link
	// with the token in the query string. When they click the link, you can validate the token and
	// set the user's email to validated in the database. I won't do this here because it's a bit
	// out of scope for this project + it would require setting up an AWS SES account or similar.
	c.Status(http.StatusCreated)
}
