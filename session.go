package main

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/sessions"
	model "github.com/joshgav/go-demo/model"
	"github.com/subosito/gotenv"
)

const (
	envvarName  = "COOKIE_KEY"
	sessionName = "vanpool_user"

	selfKey          = "self"
	authenticatedKey = "authenticated"
)

var store sessions.Store

func init() {
	gotenv.Load()
	key := os.Getenv(envvarName)
	if len(key) == 0 {
		log.Printf("Session (init): use envvar %v for cookie key\n", envvarName)
		key = "makemerandom"
	}

	log.Printf("Session (init): registering Rider type\n")
	gob.Register(&model.Rider{})
	log.Printf("Session (init): creating new cookie store\n")
	store = sessions.NewCookieStore([]byte(key))
}

// Session is middleware which creates/restores session info
// session vars: self *model.Rider, authenticated bool
// use in a later handler: // TODO: helper methods
//   `rider, ok := r.Context().Value(selfKey).(*model.Rider)`
//   `authenticated, ok := r.Context().Value(authenticatedKey).(bool)`
func Session(next http.Handler) http.Handler {
	log.Printf("Session: hello")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("SetSession: getting session %v\n", sessionName)
		s, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		rider := s.Values[selfKey].(*model.Rider)
		if rider == nil {
			log.Printf("Session: no session found\n")
		}
		s.Values[authenticatedKey] = rider == nil
		log.Printf("Session: saving session\n")
		s.Save(r, w)

		authenticated := s.Values[authenticatedKey]
		var ctx context.Context
		ctx = context.WithValue(r.Context(), selfKey, rider)
		ctx = context.WithValue(r.Context(), authenticatedKey, authenticated)

		log.Printf("Session: done, calling next with context")
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetSession(token string, w http.ResponseWriter, r *http.Request) error {
	log.Printf("SetSession: getting session %v\n", sessionName)
	s, err := store.Get(r, sessionName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, nil)
	if err != nil {
		log.Printf("failed to parse JWT: %v\n", err)
	}
	log.Printf("SetSession: parsed jwt: %+v\n", t)
	s.Values[selfKey] = &model.Rider{}
	s.Save(r, w)
	return err
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UserHandler: checking for session user\n")
	rider := r.Context().Value(selfKey).(*model.Rider)
	if rider == nil {
		// send an error message
	}
	log.Printf("UserHandler: responding with session user: %+v\n", rider)
	json, err := json.Marshal(rider)
	if err != nil {
		log.Printf("UserHandler: failed to marshal json: %v\n", err)
	}
	w.Write(json)
}
