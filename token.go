package main

import (
	"time"

	jwt "github.com/golang-jwt/jwt/v5"
)

const issuer = "my-short-url"

func isTokenValid(token string) (*jwt.Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("mySigningKey"), nil
	})

	return t, err
}

func generateUserToken(username string, userid int) (string, error) {
	type MyCustomClaims struct {
		Username string `json:"username"`
		UserID   int    `json:"user_id"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		username,
		userid,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "user-token",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("mySigningKey"))

	if err != nil {
		return "", err
	}

	return signedToken, err

}

func generateEmailConfirmationLink(username, email string) (string, error) {
	type MyCustomClaims struct {
		Username string `json:"username"`
		jwt.RegisteredClaims
	}

	// Create claims with multiple fields populated
	claims := MyCustomClaims{
		username,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   "email-confirmation",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte("mySigningKey"))

	if err != nil {
		return "", err
	}

	return signedToken, err
}
