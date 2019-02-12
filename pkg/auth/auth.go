package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joshgav/vanpool-manager/pkg/model"
	"github.com/joshgav/vanpool-manager/pkg/session"
	"github.com/joshgav/vanpool-manager/pkg/util"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/microsoft"
)

var (
	oauth2Config *oauth2.Config

	redirectURIf        = "%v://%v/login"
	redirectURIScheme   string
	redirectURIHostname string

	clientID     string
	clientSecret string
	scopes       = []string{
		"openid",
		"email",
		"profile",
		"offline_access",
		// must specify a non-OpenID scope to get access token
		// i.e. if only OpenID scopes are used only id_token
		// and refresh_token (if offline_access is requested)
		// are returned
		"user.read",
	}
)

func init() {
	redirectURIScheme := util.GetenvOrDefault("OAUTH_LOGIN_SCHEME", "https")
	redirectURIHostname := util.GetenvOrDefault("OAUTH_LOGIN_HOSTNAME", "localhost:8080")
	clientID = os.Getenv("OAUTH_CLIENT_ID")
  if clientID == "" {
    log.Fatalln("Supply a ClientID in OAUTH_CLIENT_ID")
  }
	clientSecret = os.Getenv("OAUTH_CLIENT_SECRET")
  if clientSecret == "" {
    log.Fatalln("Supply a client secret in OAUTH_CLIENT_SECRET")
  }

	oauth2Config = &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     microsoft.AzureADEndpoint(""), // Note: AADv2
		Scopes:       scopes,
		RedirectURL:  fmt.Sprintf(redirectURIf, redirectURIScheme, redirectURIHostname),
	}
}

// Authentication is net/http middleware which checks session to see
// if current user is authenticated, and if not redirects to a login server
func Authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Authentication: checking for existing authenticated session\n")
		authenticated, ok := r.Context().Value(session.AuthenticatedKey).(bool)
		log.Printf("Authentication: authenticated?: %b\n", authenticated)
		if (ok == false || authenticated == false) {
			state := r.Context().Value(session.StateKey).(string)
			log.Printf("Authentication: using state: %v\n", state)
			authorizeURL := oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
			log.Printf("Authentication: redirecting to %s\n", authorizeURL)
			http.Redirect(w, r, authorizeURL, http.StatusFound)
			return
		} else { // authenticated == true
      log.Printf("Authentication: user is authenticated, done\n")
      next.ServeHTTP(w, r)
    }
	})
}

// AuthzCodeHandler is net/http middleware which expects to receive an authz code
// from a login server. It uses this to get an OAuth access token and Open ID id_token
// and populates the session user based on their attributes.
func AuthzCodeHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("AuthzCodeHandler: extracting code and checking state\n")
	state, ok := r.Context().Value(session.StateKey).(string)
  if ok == false {
		http.Error(w, "AuthzCodeHandler: could not find state\n",
			http.StatusInternalServerError)
		return
	}
	if state != r.FormValue("state") {
		log.Printf("AuthzCodeHandler: state mismatch: have: %s; got: %s\n",
			state, r.FormValue("state"))
		http.Error(w, "AuthzCodeHandler: state doesn't match session's state, rejecting",
			http.StatusNotAcceptable)
		return
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
	log.Printf("AuthzCodeHandler: got token: %+v\n", token)

	idToken, ok := token.Extra("id_token").(string)
	if ok == false {
		log.Printf("AuthzCodeHandler: but didn't find id_token\n")
		http.Error(w, "didn't receive id_token", http.StatusInternalServerError)
		return
	} else {
		log.Printf("AuthzCodeHandler: and id_token: %+v\n\n", idToken)
	}

	log.Printf("AuthzCodeHandler: building rider via id_token: %+v\n", idToken)
	rider, err := riderFromJwt(idToken)
	if err != nil {
		log.Printf("AuthzCodeHandler: failed to build rider from jwt: %s\n", err)
		http.Error(w, "failed to build rider from jwt", http.StatusInternalServerError)
		return
	}

	log.Printf("AuthzCodeHandler: setting state with rider: %v\n", rider)
	session.SetSession(rider, w, r)
	log.Printf("AuthzCodeHandler: done, redirecting to SPA\n")
	http.Redirect(w, r, "/web/", http.StatusFound)
}

func azureKeyFunc (token *jwt.Token) (interface{}, error) {
		/*
			// retrieve from https://login.microsoftonline.com/common/v2.0/.well-known/openid-configuration
			// and https://login.microsoftonline.com/common/discovery/v2.0/keys
			// but doesn't work properly at the moment

			kid := token.Header["kid"].(string)
			// get key from discovery document, then parse and return
			verifyKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(azurePubKey))
			if err != nil {
				log.Printf("could not parse public key from string: %#v\n", err)
			}
			return verifyKey, nil
		*/
		return nil, nil // obviously not acceptable
}

// riderFromJwt takes info from an id_token to create a rider with defaults
func riderFromJwt(_jwt string) (*model.Rider, error) {
	idToken, err := jwt.Parse(_jwt, azureKeyFunc)
	if err != nil {
		log.Printf("ridersFromJwt: could not parse id_token %#v\n", err.Error())
		if err.Error() == "key is of invalid type" {
			log.Printf("ridersFromJwt: continuing despite error\n")
		} else {
			return nil, errors.New(fmt.Sprintf("could not parse id_token: %#v", err.Error()))
		}
	}
	log.Printf("riderFromJwt: parsed id_token: %#v\n", idToken)

  // jwt.MapClaims parses claims to a map[string]interface{}
	claims, ok := idToken.Claims.(jwt.MapClaims)
  if ok == false {
		return nil, errors.New("failed to parse claims in id_token\n")
	}
	rider := &model.Rider{
		ID:          uuid.Must(uuid.NewV4()),
		Username:    claims["email"].(string),
		DisplayName: claims["name"].(string),
		// defaults: now, in
		Date:      util.ParseDate(time.Now()),
		Direction: model.TravelDirectionInbound,
	}
	log.Printf("riderFromJwt: created rider %#v\n", rider)
	return rider, nil
}
