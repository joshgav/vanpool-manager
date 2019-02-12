package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"

  "github.com/joshgav/vanpool-manager/pkg/auth"
  "github.com/joshgav/vanpool-manager/pkg/handlers"
  "github.com/joshgav/vanpool-manager/pkg/session"
)

var (
	webdir = "./web/dist"
)

func configureRouter() http.Handler {
	r := mux.NewRouter()

  // redirect / to /web
	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
		return
	})

	// serve web page
	root := r.PathPrefix("/web").Subrouter()
	root.Use(session.Session)
	root.Use(auth.Authentication)
	root.PathPrefix("/").Handler(
		http.StripPrefix("/web/", http.FileServer(http.Dir(webdir))))

	// redirect unidentified users to log in
	login := r.Path("/login").Subrouter()
	login.Use(session.Session)
	login.Methods("GET").HandlerFunc(auth.AuthzCodeHandler)

  // this app's API
	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(session.Session)
	api.Use(auth.Authentication)

	// GET /api/v1/user
	api.Path("/user").Methods("GET").
		HandlerFunc(session.UserHandler)

	// GET /api/v1/riders?date=2018-12-01&direction=in
	api.Path("/riders").Methods("GET").
		Queries("date", "{date}", "direction", "{direction}").
		HandlerFunc(handlers.RidersGetHandler)

	// PUT /api/v1/riders json:*model.Rider
	api.Path("/riders").Methods("PUT").
		HandlerFunc(handlers.RidersPutHandler)

	// POST /api/v1/riders/delete json:*model.Rider
	api.Path("/riders/delete").Methods("POST").
		HandlerFunc(handlers.RidersDeleteHandler)

  return r
}

func main() {
  r := configureRouter()

	port := os.Getenv("PORT")
  if port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port) // prepend colon

	log.Printf("starting http server on port %v\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
