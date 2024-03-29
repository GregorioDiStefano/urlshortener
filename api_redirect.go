package main

import (
	"net/http"
	"strings"

	gin "github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func (app *App) redirect(c *gin.Context) {
	key := strings.TrimLeft(c.Request.URL.Path, "/")

	// Minimum length of a short url is 4 characters (2 for the base64 encoded id, 2 for the nonce)
	if len(key) < 4 {
		c.Status(http.StatusNotFound)
		return
	}

	// split short url into two parts: actual id component based on row id, and nonce
	var dbID uint64
	var nonce string
	var err error
	if dbID, nonce, err = shortURLKeyToIDAndNonce(key); err != nil {
		log.WithError(err).WithField("key", key).Error("error converting key")
		c.Status(http.StatusNotFound)
		return
	}

	// update last accessed and access count in background
	go func(dbID uint64) {
		if err := app.db.UpdateAccessAndLastAccessed(dbID); err != nil {
			log.WithError(err).WithField("id", dbID).Error("error updating access count")
		}
	}(dbID)

	// check if we have the target url in cache
	if target, err := app.cache.GetURL(key); err == nil && target != "" {
		// TODO: keep cache entry "warm"; reset ttl
		log.WithField("key", key).WithField("target", target).Info("cache hit")
		c.Redirect(http.StatusMovedPermanently, target)
		return
	}

	// it's not cached, so get it from the database
	url, nonceExpected, err := app.db.GetURL(dbID)
	if err != nil {
		log.WithError(err).WithField("id", dbID).Error("error getting url")
		c.Status(http.StatusNotFound)
		return
	}

	// check if the nonce matches
	if nonce != nonceExpected {
		log.WithField("id", dbID).WithField("nonce", nonce).WithField("expected", nonceExpected).Error("nonce mismatch")
		c.Status(http.StatusNotFound)
		return
	}

	// update cache
	if _, err := app.cache.InsertURL(key, url); err != nil {
		log.WithError(err).WithField("key", key).WithField("url", url).Error("error caching url")
		// carry on
	}

	c.Redirect(http.StatusMovedPermanently, url)
}
