package auth

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	authorizeURLFormatString = "https://%v/%v/oauth2/authorize?response_type=code&client_id=%v&redirect_uri=%v&state=%v&resource=%v"
	authorityHost            = "login.microsoftonline.com"
	tenantID                 = ""
	clientID                 = ""
	redirectURI              = ""
)

func CheckAuthenticated(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if mux.Vars()["username"] == "" {
			http.Redirect("/auth/login")
			return
		}
		h.ServeHTTP()
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// if state.Authenticated == true { follow original redirect }
	// else { redirect to OAuth code endpoint with appropriate params }
	state = ""
	resource = ""
	authorizeURL := fmt.Sprintf(
		authorizeURLFormatString,
		tenantID,
		clientID,
		redirectURI,
		state,
		resource,
	)
	http.Redirect(w, r, authorizeURL, 301)
}

func AuthzCodeHandler(w http.ResponseWriter, r *http.Request) {
	// get code from URL and submit to token endpoint, get back token
	// update session state (Authenticated = true, Username ? )
	// more session state could be filled in through next redirect (I think)
	// finally: return to originally requested URL

	// NewServicePrincipalTokenFromAuthorizationCode(oauthConfig OAuthConfig, clientID string, clientSecret string, authorizationCode string, redirectURI string, resource string, callbacks ...TokenRefreshCallback) (*ServicePrincipalToken, error) {})

	// sessions.Set()
	// redirect to /api/v1/user ? or populate here ?
}

func UserHandler(w http.ResponseWriter, r *http.Request) {

}
