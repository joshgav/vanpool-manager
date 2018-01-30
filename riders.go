package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/joshgav/go-demo/model"
)

func ridersGetHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ridersGetHandler: hello\n")

	date := r.FormValue("date")
	if len(date) < 10 { // yyyy-mm-dd
		http.Error(w, "date must be specified in query parameter", http.StatusBadRequest)
		return
	}

	direction := convertDirection(r.FormValue("direction"))
	if len(direction) < 1 {
		http.Error(w, "direction must be specified in query parameter", http.StatusBadRequest)
		return
	}

	log.Printf("ridersGetHandler: getting riders for date: %s, direction: %s\n",
		date, direction)
	fullDate, err := time.Parse("2006-01-02", date)
	riders, err := model.GetRidersFor(fullDate, model.TravelDirection(direction))
	if err != nil {
		http.Error(w, "could not get riders", http.StatusInternalServerError)
	}

	log.Printf("ridersGetHandler: returning riders: %v\n", riders)
	j, _ := json.Marshal(riders)
	_, err = w.Write(j)
	if err != nil {
		http.Error(w, "could not marshal riders into json", http.StatusInternalServerError)
	}
	log.Printf("ridersGetHandler: done\n")
	return
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

func convertDirection(dir string) string {
	if dir == "Outbound" || dir == "O" {
		return "O"
	} else {
		return "I"
	}
}
