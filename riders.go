package main

import (
	"net/http"
)

func ridersGetHandler(w http.ResponseWriter, r *http.Request) {
	// read date requested from query string
	// call into model.getRidersFor(date time.Time, direction TravelDirection) using that date
	// return a json array of *model.Rider
}

func ridersPutHandler(w http.ResponseWriter, r *http.Request) {
	// build up model.Rider from JSON
	// call model.AddRider(rider)
}

func ridersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	// build up model.Rider from JSON
	// call model.DeleteRider(rider)
	// match must be exact { name, date, direction }
}
