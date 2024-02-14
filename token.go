package main

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

const issuer = "my-short-url"

// jwtAuthMiddleware is a middleware that checks for a valid JWT token
// in the Authorization header and sets the username and user_id in the
// context if the token is valid
type jwtAuthMiddleware struct {
	secret []byte
}

func (middleware jwtAuthMiddleware) GenerateUserToken(userid int) (string, error) {
	type MyCustomClaims struct {
		UserID int `json:"user_id"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		userid,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(3 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user-token",
			Issuer:    issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(middleware.secret)

	if err != nil {
		return "", err
	}

	return signedToken, err
}

func (middleware jwtAuthMiddleware) ValidateToken(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return middleware.secret, nil
	})

	return t, err
}

func (middleware jwtAuthMiddleware) JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer") {
			c.JSON(401, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		actualToken := strings.TrimPrefix(authHeader, "Bearer ")
		actualToken = strings.TrimSpace(actualToken)

		if actualToken == "" {
			c.JSON(401, gin.H{"error": "token is required"})
			c.Abort()
			return
		}

		if t, err := middleware.ValidateToken(actualToken); err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		} else if t.Valid {
			var user_id int

			// Safe type assertion for username
			if claims, ok := t.Claims.(jwt.MapClaims); ok {
				if userIDClaim, ok := claims["user_id"].(float64); ok {
					user_id = int(userIDClaim)
				} else {
					c.JSON(401, gin.H{"error": "invalid token"})
					c.Abort()
					return
				}
			} else {
				c.JSON(401, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}

			if user_id != 0 {
				c.Set("user_id", user_id)
				c.Next()
			} else {
				c.JSON(401, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		}
	}
}
