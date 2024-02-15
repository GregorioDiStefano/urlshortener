package main

import (
	"fmt"
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

func (middleware jwtAuthMiddleware) GenerateUserToken(userid uint64) (string, error) {
	type LoginClaims struct {
		UserID uint64 `json:"user_id"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := LoginClaims{
		userid,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
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

		token := strings.TrimPrefix(authHeader, "Bearer ")
		token = strings.TrimSpace(token)

		if token == "" {
			c.JSON(401, gin.H{"error": "token is required"})
			c.Abort()
			return
		}

		if t, err := middleware.ValidateToken(token); err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		} else if t.Valid {
			var userID uint64

			// Safe type assertion for username
			if claims, ok := t.Claims.(jwt.MapClaims); ok {
				if userIDClaim, ok := claims["user_id"].(float64); ok {
					userID = uint64(userIDClaim)
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

			if userID != 0 {
				fmt.Println("setting user_id: ", userID)
				c.Set("user_id", userID)
				c.Next()
			} else {
				c.JSON(401, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		}
	}
}
