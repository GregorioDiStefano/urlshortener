package main

import (
	"net/http"

	gin "github.com/gin-gonic/gin"
)

// ping is essentially a health check, which would be used by k8s, etc
// to determine if the service is healthy
func (app *App) ping(c *gin.Context) {
	errCache := make(chan error)
	errDB := make(chan error)

	// run go routine
	go func() {
		errCache <- app.cache.Ping()
		errDB <- app.db.Ping()
	}()

	// wait for channels to respond with ping status
	errCacheResp := <-errCache
	errDBResp := <-errDB

	if errCacheResp == nil && errDBResp == nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
		return
	}

	c.JSON(http.StatusInternalServerError, gin.H{
		"message": "error",
	})
}
