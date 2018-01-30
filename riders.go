package main

import (
	"encoding/json"
	"io/ioutil"
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
	log.Printf("ridersPutHandler: hello")
	_json, err := ioutil.ReadAll(r.Body)
	log.Printf("ridersPutHandler: body: %s\n", _json)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		log.Printf("failed to read request body: %v\n", err)
		return
	}
	var rider *model.Rider
	err = json.Unmarshal(_json, &rider)
	if err != nil {
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		log.Printf("failed to unmarshal json: %v\n", err)
		return
	}
	log.Printf("ridersPutHandler: adding rider: %+v\n", rider)
	err = model.AddRider(rider)
	if err != nil {
		http.Error(w, "could not persist rider", http.StatusInternalServerError)
		log.Printf("failed to persist rider: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	return
}

func ridersDeleteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("ridersDeleteHandler: hello")
	_json, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "could not read request body", http.StatusBadRequest)
		log.Printf("failed to read request body: %v\n", err)
		return
	}
	var rider *model.Rider
	err = json.Unmarshal(_json, &rider)
	if err != nil {
		http.Error(w, "could not unmarshal json", http.StatusBadRequest)
		log.Printf("failed to unmarshal json: %v\n", err)
		return
	}
	log.Printf("ridersDeleteHandler: deleting rider: %+v\n", rider)
	err = model.DeleteRider(rider)
	if err != nil {
		http.Error(w, "could not delete rider", http.StatusInternalServerError)
		log.Printf("failed to delete rider: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusAccepted)
	return
}

func convertDirection(dir string) string {
	if dir == "Outbound" || dir == "O" {
		return "O"
	} else {
		return "I"
	}
}
