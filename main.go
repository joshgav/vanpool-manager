package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

var (
	webdir = "./web/dist/"
	port   = GetenvOrDefault("PORT", "8080")
)

func main() {
	r := mux.NewRouter()

	// client app (SPA)
	root := r.PathPrefix("/web").Subrouter()
	root.Use(Session)
	root.Use(Authentication)
	root.PathPrefix("/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir(webdir))))

	// OAuth authorization code handler
	login := r.Path("/login").Subrouter()
	login.Use(Session)
	login.Methods("GET").HandlerFunc(AuthzCodeHandler)

	api := r.PathPrefix("/api/v1").Subrouter()
	api.Use(Session)
	api.Use(Authentication)

	// GET /api/v1/riders?date=2018-01-05&direction=in
	api.Path("/riders").Methods("GET").
		Queries("date", "{date}", "direction", "{direction}").
		HandlerFunc(ridersGetHandler)

	// PUT /api/v1/riders json:*model.Rider
	api.Path("/riders").Methods("PUT").
		HandlerFunc(ridersPutHandler)

	// DELETE /api/v1/riders json:*model.Rider
	api.Path("/riders").Methods("DELETE").
		HandlerFunc(ridersDeleteHandler)

	api.Path("/user").Methods("GET").
		HandlerFunc(UserHandler)

	log.Printf("starting http server on port %v\n", port)
	http.ListenAndServe(":"+port, r)
}
