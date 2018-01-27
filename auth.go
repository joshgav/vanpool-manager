package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	redirectURIf    = "http://%v/login"
	redirectURIHost string
	scope           = "openid"
	tenantID        string
	clientID        string
	clientSecret    string
	oauth2Config    *oauth2.Config
)

func init() {
	hostname := os.Getenv("WEBAPP_HOSTNAME")
	if len(hostname) > 0 {
		redirectURIHost = hostname
	} else {
		redirectURIHost = "localhost:8080"
	}
	tenantID = os.Getenv("AZ_TENANT_ID") // would "common" be appropriate for a multi-tenant app?
	clientID = os.Getenv("AZ_CLIENT_ID")
	clientSecret = os.Getenv("AZ_CLIENT_SECRET")

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint(""),
		Scopes:       []string{"openid", "email", "profile"},
		RedirectURL:  fmt.Sprintf(redirectURIf, redirectURIHost),
	}
}

func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authentication: checking for existing authenticated session\n")
		var authenticated bool = false
		authenticated, _ = r.Context().Value(authenticatedKey).(bool)
		log.Printf("Authentication: authenticated: %s\n", authenticated)
		if authenticated == false {
			var state, _ = r.Context().Value(stateKey).(string)
			authorizeURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
			// should be done with net/url
			// authorizeURL = strings.Join([]string{authorizeURL, "&response_mode=form_post"}, "")
			log.Printf("Authentication: redirecting to %s\n", authorizeURL)
			http.Redirect(w, r, authorizeURL, 301)
		}
		log.Printf("Authentication: user is authenticated, done\n")
		next.ServeHTTP(w, r)
	})
}

func AuthzCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthzCodeHandler: extracting code and checking state\n")
	var ok bool
	var state string
	if state, ok = r.Context().Value(stateKey).(string); ok == false {
		// http.Error("there has to be state")
	}
	if state != r.FormValue("state") {
		log.Printf("state mismatch: state: %s; r.FormValue(\"state\"): %s\n", state, r.FormValue("state"))
		// http.Redirect("/")
	}

	code := r.FormValue("code")
	log.Printf("AuthzCodeHandler: going to request access token with code: %s\n", code)
	token, err := oauth2Config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalf("AuthzCodeHandler: failed to get access token with authz code: %v\n", err)
		http.Error(w, "failed to get access token with authz code", http.StatusInternalServerError)
		http.Redirect(w, r, "/error", 301)
	}
	log.Printf("AuthzCodeHandler: setting session via token: %+v\n", token)
	SetSession(token.AccessToken, w, r)
	log.Printf("AuthzCodeHandler: done, redirecting to SPA\n")
	http.Redirect(w, r, "/", 301)
}
