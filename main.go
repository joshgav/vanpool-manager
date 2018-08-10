package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

var (
	webdir = "./web/dist"
)

func main() {
	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port) // prepend colon

	r := mux.NewRouter()

	r.Path("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/web/", http.StatusMovedPermanently)
		return
	})

	// client app (SPA)
	root := r.PathPrefix("/web").Subrouter()
	root.Use(Session)
	root.Use(Authentication)
	root.PathPrefix("/").Handler(
		http.StripPrefix("/web/", http.FileServer(http.Dir(webdir))))

	// OAuth authorization code handler
	login := r.Path("/login").Subrouter()
	login.Use(Session)
	login.Methods("GET").HandlerFunc(AuthzCodeHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(Session)
	api.Use(Authentication)

	// GET /api/v1/user
	api.Path("/user").Methods("GET").
		HandlerFunc(UserHandler)

	// GET /api/v1/riders?date=2018-12-01&direction=in
	api.Path("/riders").Methods("GET").
		Queries("date", "{date}", "direction", "{direction}").
		HandlerFunc(ridersGetHandler)

	// PUT /api/v1/riders json:*model.Rider
	api.Path("/riders").Methods("PUT").
		HandlerFunc(ridersPutHandler)

	// POST /api/v1/riders/delete json:*model.Rider
	api.Path("/riders/delete").Methods("POST").
		HandlerFunc(ridersDeleteHandler)

	log.Printf("starting http server on port %v\n", port)
	log.Fatal(http.ListenAndServe(port, r))
}
