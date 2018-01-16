package main

import (
	"github.com/gin-gonic/gin"
)

func main(args []string) {
	r := gin.Default()
	r.Use(sessionStateHandler())
	r.Use(authenticationHandler())

	r.Static("/", "./dist/assets")

	auth := r.Group("/auth")
	auth.GET("/login", authLoginHandler())
	auth.POST("/code", authCodeHandler())

	api := r.Group("/api/v1")

	api.GET("/riders", ridersGetHandler())
	api.PUT("/riders", ridersPutHandler())
	api.DELETE("/riders", ridersDeleteHandler())

	api.GET("/currentUser", currentUserHandler())
}
