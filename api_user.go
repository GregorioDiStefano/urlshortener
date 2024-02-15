package main

import (
	"fmt"
	"net/http"

	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// disableURL is a handler function that takes a short URL and disables it.
// NOTE: we purposely don't delete the URL from the database, as we want to keep a record of it
func (app *App) disableURL(c *gin.Context) {
	key := c.Param("id")
	userID := c.GetUint64("user_id")

	var dbID uint64
	var err error

	if dbID, _, err = shortURLKeyToIDAndNonce(key); err != nil {
		log.WithError(err).WithField("key", key).Error("error converting key")
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid key"})
		return
	}

	if err := app.db.DisableURL(userID, dbID); err != nil {
		log.WithField("user_id", userID).WithField("key", key).WithError(err).Error("error disabling url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := app.cache.DeleteURL(key); err != nil {
		log.WithField("key", key).WithError(err).Error("error deleting url from cache")
	}

	c.Status(http.StatusOK)
}

func (app *App) urls(c *gin.Context) {
	userID := c.GetUint64("user_id")
	urls, err := app.db.GetURLs(userID)

	if err != nil {
		log.WithField("user_id", userID).WithError(err).Error("error getting urls")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, urls)
}

// shortenURL is a handler function that takes a URL, shortens it, and returns the shortened URL.

// NOTE: the shorten URL key is essentially the base64 encoding of the row id in the database, and the nonce
// this is done to make enumeration of short urls more difficult
func (app *App) shortenURL(c *gin.Context) {
	type ShortenRequest struct {
		URL string `json:"url" binding:"required,http_url"`
	}

	userId := c.GetUint64("user_id")
	var json ShortenRequest
	if err := c.ShouldBindJSON(&json); err != nil {
		log.WithError(err).Error("error binding json")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, nonce, err := app.db.InsertURL(userId, json.URL)
	if err != nil {
		log.WithField("user_id", userId).WithError(err).Error("error inserting url")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	key := uint64ToBase64(userID)
	shortURL := fmt.Sprintf("%s%s", key, nonce)
	c.JSON(http.StatusCreated, gin.H{"key": shortURL})
}
