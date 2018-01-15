package main

const (
	authFString = "https://login.microsoftonline.com/%v/oauth2/authorize?response_type=code&client_id=%v&redirect_uri=%v&state=%v&resource=%v"
	tenantID    = ""
	clientID    = ""
	redirectURI = ""
)

func authenticationHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if state.Authenticated == true { continue to original URL }
		// if path doesn't include `/auth` { redirect to /login }
		// if it does include /auth don't redirect, we might be going to /code
		http.Redirect("/auth/login")
		c.Next()
	}
}

func authLoginHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// if state.Authenticated == true { ?? }
		// else { redirect to OAuth code endpoint with appropriate params }
		authorizationURL := fmt.Sprintf(
			authFSString,
			tenantID,
			clientID,
			redirectURI,
			stateID, // session state key
			resource,
		)
		// update for gin:
		// http.Redirect(w, r, authorizationURL, 301)
		c.Next()
	}
}

func authCodeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get code from URL and submit to token endpoint, get back token
		// update session state (Authenticated = true, Username ? )
		// more session state could be filled in through next redirect (I think)
		// finally: return to originally requested URL

		//NewServicePrincipalTokenFromAuthorizationCode(oauthConfig OAuthConfig, clientID string, clientSecret string, authorizationCode string, redirectURI string, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {})
	}
}
