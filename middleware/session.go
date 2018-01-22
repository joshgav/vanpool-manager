package middleware

import (
	"encoding/gob"
	"fmt"
	"net/http"
	"strings"

	"github.com/joshgav/go-demo/riders"

	"github.com/gorilla/sessions"
	"github.com/subosito/gotenv"
)

const (
	envvarName  = "COOKIE_KEY"
	sessionName = "vanpool_user"

	SelfKey          = "self"
	AuthenticatedKey = "authenticated"
)

// Session is middleware which restores session info
// session vars: self *riders.Rider, authenticated bool
// use in a later handler:
//   `rider, ok := r.Context().Value(middleware.SelfKey).(*riders.Rider)`
//   `authenticated, ok := r.Context().Value(middleware.AuthenticatedKey).(bool)`
func Session(next http.Handler) http.Handler {

	key, _ := os.Getenv(envvarName)
	if len(strings.Trim(key)) == 0 {
		panic(fmt.Sprintf("specify envvar %s for sessions", envvarName))
	}

	gob.Register(&riders.Rider{})
	var store = sessions.NewCookieStore([]byte(key))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s, err := store.Get(r, sessionName)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		rider, ok := session.Values[selfKey].(*riders.Rider)
		if !ok {
			http.Redirect("/auth/login", 301)
		}
		session.Values[authenticatedKey] = rider != nil
		authenticated := session.Values[authenticatedKey]
		session.Save(r, w)

		ctx := context.WithValue(r.Context(), SelfKey, rider)
		ctx = context.WithValue(r.Context(), AuthenticatedKey, authenticated)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
