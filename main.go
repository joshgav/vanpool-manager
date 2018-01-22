package main

import (
	"net/http"

	mux "github.com/gorilla/mux"

	auth "github.com/joshgav/go-demo/auth"
	mw "github.com/joshgav/go-demo/middleware"
	riders "github.com/joshgav/go-demo/riders"
)

const dir = "./web"

func main(args []string) {
	r := mux.NewRouter()

	// GET /
	r.Path("/").
		Use(mw.Session).
		Use(auth.CheckAuthenticated).
		Handler(http.StripPrefix("/web/", http.FileServer(http.Dir(dir))))

	s := r.PathPrefix("/auth").Subrouter
	// GET /auth/login
	s.Path("/login").Methods("GET").Handler(auth.LoginHandler)
	// POST /auth/callback
	s.Path("/callback").Methods("POST").Handler(auth.AuthzCodeHandler)

	api := r.PathPrefix("/api/v1").Subrouter.
		Use(mw.Session).
		Use(auth.CheckAuthenticated)

	// query parameters are coincidentally the same but shouldn't be DRYed
	// GET /api/v1/riders?date=2018-01-05&direction=in
	api.Path("/riders").Methods("GET").
		Queries("date", "{date}", "direction", "{direction}").
		Handler(riders.GetHandler)
	// PUT /api/v1/riders?date=2018-01-05&direction=in
	api.Path("/riders").Methods("PUT").
		Queries("date", "{date}", "direction", "{direction}").
		Handler(riders.PutHandler)
	// DELETE /api/v1/riders?date=2018-01-05&direction=in
	api.Path("/riders").Methods("DELETE").
		Queries("date", "{date}", "direction", "{direction}").
		Handler(riders.DeleteHandler)

	api.Path("/user").Methods("GET").
		Handler(auth.UserHandler)

	http.ListenAndServe(":"+port, r)
}
