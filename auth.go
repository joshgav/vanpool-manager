package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	// jwt "github.com/dgrijalva/jwt-go"
	"github.com/joshgav/go-demo/model"

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
		Scopes:       []string{"openid", "email", "profile", "offline_access"},
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
			log.Printf("Authentication: using state: %v\n", state)
			authorizeURL := oauth2Config.AuthCodeURL(state,
				oauth2.AccessTypeOffline)
			// add `oauth2.SetAuthURLParam("response_mode", "form_post")` to arry
			log.Printf("Authentication: redirecting to %s\n", authorizeURL)
			http.Redirect(w, r, authorizeURL, 301)
			return
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
		log.Printf("AuthzCodeHandler: failed to exchange authz code: %v\n", err)
		http.Error(w, "failed to get access token with authz code",
			http.StatusInternalServerError)
		return
	}
	log.Printf("AuthzCodeHandler: got token: %s\n", token)
	if idToken, ok := token.Extra("id_token").(string); ok == false {
		log.Printf("AuthzCodeHandler: but didn't find id_token")
	} else {
		log.Printf("AuthzCodeHandler: and id_token: %s\n", idToken)
	}

	/*
		// initial token from token endpoint contains only refresh_token
		// stash it in a new TokenSource so it will be auto-refreshed
		// unfortunately when access_token isn't set oauth2.Config.Exchange(...) only returns
		// an error and token is nil
		// once it is modified to not return an error,
		//   the initial token seems to be sufficient

		ts := oauth2Config.TokenSource(ctx, token)
		c := oauth2.NewClient(ctx, ts)
		ctx = context.WithValue(ctx, oauth2.HTTPClient, c)
		token2, err := oauth2Config.Exchange(ctx, code)
		log.Printf("AuthzCodeHandler: got another token: %v\n", token2)
		log.Printf("AuthzCodeHandler: got error: %v\n", err)
		if err != nil {
			log.Fatalf("AuthzCodeHandler: failed to get access token with refresh token: %v\n", err)
		}
	*/

	log.Printf("AuthzCodeHandler: building rider via token: %+v\n", token)
	rider, err := riderFromJwt(token.AccessToken)
	if err != nil {
		log.Printf("AuthzCodeHandler: failed to build rider from jwt: %v\n", err)
		http.Error(w, "failed to build rider from jwt", http.StatusInternalServerError)
		return
	}
	log.Printf("AuthzCodeHandler: setting state with rider: %v\n", rider)
	SetSession(rider, w, r)
	log.Printf("AuthzCodeHandler: done, redirecting to SPA\n")
	http.Redirect(w, r, "/web/", 301)
}

func riderFromJwt(_jwt string) (*model.Rider, error) {
	// jwt.Parse(_jwt)
	return &model.Rider{}, nil
}
