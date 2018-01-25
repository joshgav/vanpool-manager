package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Azure/go-autorest/autorest/adal"
	"github.com/Azure/go-autorest/autorest/azure"
)

const (
	authorizeURLf = "https://%v/%v/oauth2/authorize?response_type=code+id_token&client_id=%v&redirect_uri=%v&state=%v&scope=%v"
	authorizeHost = "login.microsoftonline.com"
	redirectURIf  = "http://%v/login"
	scope         = "openid"
)

var (
	redirectURIHost string
	state           string
	tenantID        string
	clientID        string
	clientSecret    string
)

func init() {
	hostname := os.Getenv("WEBAPP_HOSTNAME")
	if len(hostname) > 0 {
		redirectURIHost = hostname
	} else {
		redirectURIHost = "localhost:8080"
	}
	state = "makemerandom"
	tenantID = os.Getenv("AZ_TENANT_ID") // would "common" be appropriate for a multi-tenant app?
	clientID = os.Getenv("AZ_CLIENT_ID")
	clientSecret = os.Getenv("AZ_CLIENT_SECRET")
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CheckAuthenticatedHandler: checking for existing authenticated session\n")
		authenticated := r.Context().Value(authenticatedKey)
		log.Printf("CheckAuthenticatedHandler: authenticated: %s\n", authenticated.(bool))
		if authenticated == false {
			authorizeURL := fmt.Sprintf(
				authorizeURLf,
				tenantID,
				clientID,
				fmt.Sprintf(redirectURIf, redirectURIHost),
				state,
				scope,
			)
			log.Printf("CheckAuthenticatedHandler: redirecting to %s\n", authorizeURL)
			http.Redirect(w, r, authorizeURL, 301)
		}
		next.ServeHTTP(w, r)
	})
}

func AuthzCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthzCodeHandler: extracting code and checking state\n")
	code := r.FormValue("code")
	if state != r.FormValue("state") {
		log.Printf("state mismatch: state: %s; r.FormValue(\"state\"): %s\n", state, r.FormValue("state"))
		// http.Redirect("/")
	}

	log.Printf("AuthzCodeHandler: going to request access token\n")
	config, _ := adal.NewOAuthConfig(azure.PublicCloud.ActiveDirectoryEndpoint, tenantID)
	token, err := adal.NewServicePrincipalTokenFromAuthorizationCode(
		*config,
		clientID,
		clientSecret,
		code,
		fmt.Sprintf(redirectURIf, redirectURIHost),
		azure.PublicCloud.ActiveDirectoryEndpoint, // I think? maybe shouldn't pass `resource` at all, but `scope`?
		nil, // callbacks ...TokenRefreshCallback
	)

	if err != nil {
		log.Printf("failed to get access token: %v", err)
		// http.Redirect("/")
	}
	log.Printf("AuthzCodeHandler: setting session via token: %+v\n", token)
	SetSession(token.AccessToken, w, r)
	log.Printf("AuthzCodeHandler: done, redirecting to SPA\n")
	http.Redirect(w, r, "/", 301)
}
