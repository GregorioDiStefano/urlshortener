package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

// TODO : actual validation is missing in json
type App struct {
	db              Database
	cache           Cache
	tokenMiddleware TokenValidator
}

type TokenValidator interface {
	GenerateUserToken(userid int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	JWTMiddleware() gin.HandlerFunc
}

func (app *App) ping(c *gin.Context) {
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
}

func getRequiredEnvVar(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required", key)
	}

	return value
}

func main() {
	// Very primative way to handle config, but time constraints..
	// I would usually use https://github.com/spf13/viper for this
	dbConfig := &dbConfig{}
	dbConfig.host = getRequiredEnvVar("DB_HOST")
	dbConfig.port = getRequiredEnvVar("DB_PORT")
	dbConfig.user = getRequiredEnvVar("DB_USER")
	dbConfig.password = getRequiredEnvVar("DB_PASSWORD")
	dbConfig.dbname = getRequiredEnvVar("DB_NAME")

	jwtSecret := getRequiredEnvVar("JWT_SECRET")

	db, err := NewDB(dbConfig)

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

	tokenChecker := jwtAuthMiddleware{secret: []byte(jwtSecret)}

	app := &App{db, cache, tokenChecker}
	engine := setupRouter(app)

	engine.Run(":8888")
}
