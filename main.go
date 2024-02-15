package main

import (
	"os"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	log "github.com/sirupsen/logrus"
)

type App struct {
	db              Database
	cache           Cache
	tokenMiddleware TokenValidator
}

type TokenValidator interface {
	GenerateUserToken(userid uint64) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
	JWTMiddleware() gin.HandlerFunc
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

	exitCode := 0

	// cleanup and exit; required to ensure that the defer statements are executed
	defer func() {
		os.Exit(exitCode)
	}()

	// db env variables
	dbConfig := &dbConfig{}
	dbConfig.host = getRequiredEnvVar("DB_HOST")
	dbConfig.port = getRequiredEnvVar("DB_PORT")
	dbConfig.user = getRequiredEnvVar("DB_USER")
	dbConfig.password = getRequiredEnvVar("DB_PASSWORD")
	dbConfig.dbname = getRequiredEnvVar("DB_NAME")

	// jwt secret
	jwtSecret := getRequiredEnvVar("JWT_SECRET")

	// cache secret
	redisHost := getRequiredEnvVar("REDIS_HOST")

	db, err := NewDB(dbConfig)

	if err != nil {
		log.Warn(err)
		exitCode = 1
		return
	}

	defer db.Close()
	cache, err := NewCache(redisHost)

	if err != nil {
		log.Warn(err)
		exitCode = 1
		return
	}

	defer cache.Close()

	log.SetFormatter(&log.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	log.Info("Starting server")

	tokenChecker := jwtAuthMiddleware{secret: []byte(jwtSecret)}

	app := &App{db, cache, tokenChecker}
	engine := setupRouter(app)

	if err := engine.Run(":8888"); err != nil {
		log.Warn("Server failed to start: ", err)
		exitCode = 1
	}
}
