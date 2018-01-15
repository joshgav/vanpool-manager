package main

import (
	"github.com/gin-gonic/gin"
)

func sessionStateHandler() gin.Handler {
	return func(c *gin.Context) {
		// check for session cookie
		// if cookie is present, use to rehydrate state from cache
		// check that state is all there, including user name and profile data
		//	 if not call directory to populate
		// if no cookie create a new SessionState object, persist to cache, set cookie
		//   initial state should have Authenticated == false
	}
}
