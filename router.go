package main

import "github.com/gin-gonic/gin"

func setupRouter(app *App) *gin.Engine {
	router := gin.Default()

	/* Main router used to redirect from shortURL */
	router.GET("/:id", app.redirect)

	/* Routes used to authenticate */
	auth := router.Group("/auth")
	auth.POST("/login", app.login)
	auth.POST("/register", app.register)

	/* Authenticated Endpoints */
	api := router.Group("/api/v1")
	api.Use(app.tokenMiddleware.JWTMiddleware())
	api.POST("/shorten", app.shortenURL)
	api.DELETE("/shorten/:id", app.disableURL)
	api.GET("/urls/", app.urls)

	/* Health check for k8s; respond with http 200 if redis and postgresql are up */
	router.GET("/ping", app.ping)
	return router
}
