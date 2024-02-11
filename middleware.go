package main

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
)

func jwtAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
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

		if t, err := isTokenValid(actualToken); err != nil {
			c.JSON(401, gin.H{"error": "invalid token"})
			c.Abort()
			return
		} else if t.Valid {
			var username string
			var id int

			// Safe type assertion for username
			if claims, ok := t.Claims.(jwt.MapClaims); ok {
				if usernameClaim, ok := claims["username"].(string); ok {
					username = usernameClaim
				} else {
					c.JSON(401, gin.H{"error": "invalid token"})
					c.Abort()
					return
				}

				if userIDClaim, ok := claims["user_id"].(float64); ok {
					id = int(userIDClaim)
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

			if username != "" {
				c.Set("username", username)
				c.Set("user_id", id)
				fmt.Println("username", username)
				c.Next()
			} else {
				c.JSON(401, gin.H{"error": "invalid token"})
				c.Abort()
				return
			}
		}
	}
}
