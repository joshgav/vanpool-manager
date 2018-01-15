package main

import (
	"github.com/gin-gonic/gin"
)

func ridersGetHandler() *gin.Handler {
	return func(c *gin.Context) {
		// read date requested from query string
		// call into data.getRidersFor5Days(date time.Time) using that date
		// return a json array of riders
	}
}

func ridersPutHandler() *gin.Handler {
	return func(c *gin.Context) {
		// build up data.Rider from JSON
		// call data.AddRider(data.Rider)
	}
}

func ridersDeleteHandler() *gin.Handler {
	return func(c *gin.Context) {
		// call data.DeleteRider(data.Rider)
		// match must be exact { name, date, direction }
	}
}
