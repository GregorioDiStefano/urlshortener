package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// TODO : actual validation is missing in json
type App struct {
	db    Database
	cache Cache
}

type RestError struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Error   string `json:"error"`
}

func (app *App) validate(c *gin.Context) {
	tokenFromQuery := c.Query("token")
	if tokenFromQuery == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token is required"})
		return
	}

	t, err := isTokenValid(tokenFromQuery)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	if t.Valid {
		username := t.Claims.(jwt.MapClaims)["username"].(string)
		if err := app.db.ValidateUser(username); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to validate"})
		} else {
			c.JSON(http.StatusNoContent, nil)
		}
	}
}

func (app *App) login(c *gin.Context) {
	type Login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var json Login
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user *User
	var err error
	if user, err = app.db.GetUser(json.Username); err != nil {
		fmt.Println(user, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword(user.password_hash, []byte(json.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := generateUserToken(user.username, user.id)
	if err != nil {
		log.Warn("error generating token: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (app *App) register(c *gin.Context) {
	type RegisterRequest struct {
		Username string `json:"username" binding:"required,alphanum"`
		Password string `json:"password" binding:"required,ascii"`
		Email    string `json:"email" binding:"required,email"`
	}

	validate := func(login *RegisterRequest) error {
		const minPasswordLength = 8
		const minUsernameLength = 5

		const maxUsernameLength = 255
		const maxPasswordLength = 255

		const specialChars = "!@#$%^&*()_+"

		// TODO: more validation needed here
		if len(login.Username) < minUsernameLength {
			return fmt.Errorf("username must be at least %d characters long and at most:", minUsernameLength, maxUsernameLength)
		}
		if len(login.Password) < minPasswordLength {
			return fmt.Errorf("password must be at least %d characters long and at most", minPasswordLength, maxPasswordLength)
		}
		if !strings.ContainsAny(login.Password, specialChars) {
			return fmt.Errorf("password must contain at least one special character")
		}

		return nil
	}

	var json RegisterRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validate(&json); err != nil {
		// TODO: use a better error message
		c.JSON(http.StatusBadRequest, gin.H{"error": "validation failed: " + err.Error()})
		return
	}

	if err := app.db.SignupUser(json.Username, json.Password, json.Email); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	c.Status(http.StatusCreated)
	// Usually you would send a validation email here, i have the code for that, but I won't do anything
	// with it, since it's not the focus of the exercise
}

func (app *App) shortenURL(c *gin.Context) {
	type ShortenRequest struct {
		URL string `json:"url" binding:"required,http_url"`
	}

	var json ShortenRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := c.GetInt("user_id")
	_, _, err := app.db.InsertURL(userId, json.URL)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

}

func (app *App) shortenURLDelete(c *gin.Context) {
}

func (app *App) redirect(c *gin.Context) {
	key := strings.TrimLeft(c.Request.URL.Path, "/")
	dbID, err := keyToID(key)

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	url, err := app.db.GetURL(dbID)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusMovedPermanently, url)
}

func (app *App) urls(c *gin.Context) {
	userId := c.GetInt("user_id")
	urls, err := app.db.GetURLs(userId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	fmt.Println(urls)

	c.JSON(http.StatusOK, urls)
}

func setupRouter(app *App) *gin.Engine {
	router := gin.Default()

	router.GET("/:id", app.redirect)

	auth := router.Group("/auth")
	auth.POST("/login", app.login)
	auth.POST("/register", app.register)
	auth.GET("/validate", app.validate)

	/* Authenticated End points */
	api := router.Group("/api/v1")
	api.Use(jwtAuthMiddleware())
	api.POST("/shorten", app.shortenURL)
	api.DELETE("/shorten/:id", app.shortenURLDelete)
	api.GET("/urls/", app.urls)

	// health check for k8s
	router.GET("/ping", func(c *gin.Context) {
		errCache := make(chan error)
		errDB := make(chan error)

		go func() {
			errCache <- app.cache.Ping()
			errDB <- app.db.Ping()
		}()

		errCacheResp := <-errCache
		errDBResp := <-errDB

		if errCacheResp == nil && errDBResp == nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
			return
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "error",
			})
		}
	})

	return router
}

func main() {
	db, err := NewDB()

	if err != nil {
		log.Fatal(err)
	}

	cache, err := NewCache()

	if err != nil {
		log.Fatal(err)
	}

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.Info("Starting server")

	app := &App{db, cache}
	engine := setupRouter(app)

	engine.Run(":8888")
}