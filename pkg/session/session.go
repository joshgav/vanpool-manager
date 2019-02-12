package session

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joshgav/vanpool-manager/pkg/model"

	uuid "github.com/satori/go.uuid"
	"github.com/gorilla/sessions"
	"github.com/subosito/gotenv"
)

const (
	sessionKeyEnvVar = "SESSION_KEY"
	sessionName      = "vanpool_user"

  // used in sessions
  SelfKey = "self"
  AuthenticatedKey = "authenticated"
  StateKey = "state"
)

var store sessions.Store

func init() {
  // load env vars from local .env
	gotenv.Load()

	sessionKey := os.Getenv(sessionKeyEnvVar)
  if len(sessionKey) == 0 {
		log.Printf("Session (init): use envvar %v for a session key\n", sessionKeyEnvVar)
		sessionKey = "makemerandom"
	}

	// serialization for cookies uses encoding/gob so
  // we must register any complex types to persist
  gob.Register(&model.Rider{})
	log.Printf("Session (init): creating new cookie store\n")
	store = sessions.NewCookieStore([]byte(sessionKey))
}

// Session is middleware which creates/restores session info
// session vars: state string, self *model.Rider, authenticated bool
// use in a later handler: // TODO: helper methods
//   `rider, ok := r.Context().Value(selfKey).(*model.Rider)`
//   `authenticated, ok := r.Context().Value(authenticatedKey).(bool)`
//   `state, ok := r.Context().Value(stateKey).(string)`
func Session(next http.Handler) http.Handler {
	log.Printf("Session: hello")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Session: getting session %v\n", sessionName)
		s, err := store.Get(r, sessionName)
    if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		if _, ok := s.Values[StateKey].(string); ok == false {
			log.Printf("Session: no state currently in session, adding it\n")
			s.Values[StateKey] = uuid.Must(uuid.NewV4())
			log.Printf("Session: state set to: %v\n", s.Values[StateKey])
		}

		if _, ok := s.Values[AuthenticatedKey].(bool); ok == false {
			log.Printf("Session: user not previously authenticated\n")
			log.Printf("Session: marking user as not authenticated\n")
			s.Values[AuthenticatedKey] = false
		}

		if _, ok := s.Values[SelfKey].(*model.Rider); ok == false {
			log.Printf("Session: user not available in session\n")
      log.Printf("Session: resetting\n")
			s.Values[SelfKey] = &model.Rider{}
			s.Values[AuthenticatedKey] = false
		} else {
			log.Printf("Session: found user in session, marking as authenticated\n")
			s.Values[AuthenticatedKey] = true
		}

		log.Printf("Session: saving session (%v)\n", s)
		_ = s.Save(r, w)

		log.Printf("Session: adding session data to context for later handlers\n")
		var ctx = r.Context()
		ctx = context.WithValue(ctx, SelfKey, s.Values[SelfKey])
		ctx = context.WithValue(ctx, AuthenticatedKey, s.Values[AuthenticatedKey])
		ctx = context.WithValue(ctx, StateKey, s.Values[StateKey])

		log.Printf("Session: done, calling next with context (%v)\n", ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func SetSession(rider *model.Rider, w http.ResponseWriter, r *http.Request) error {
	log.Printf("SetSession: getting session %v\n", sessionName)
	s, err := store.Get(r, sessionName)
  if err != nil {
		http.Error(w, "failed to get session", http.StatusInternalServerError)
		return err
	}

	s.Values[SelfKey] = rider
	err = s.Save(r, w)
	return err
}

func UserHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("UserHandler: checking for session user\n")
	rider, ok := r.Context().Value(SelfKey).(*model.Rider)
	if ok == false {
		log.Printf("UserHandler: no session user found\n")
		http.Error(w, "No session user found.", http.StatusInternalServerError)
	}
	log.Printf("UserHandler: responding with session user: %+v\n", rider)
	json, err := json.Marshal(rider)
  if err != nil {
		log.Printf("UserHandler: failed to marshal json: %v\n", err)
		http.Error(w, "Failed to marshal JSON for user.", http.StatusInternalServerError)
	}
	w.Write(json)
}
