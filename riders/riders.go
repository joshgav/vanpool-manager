package riders

import (
	"github.com/gorilla/mux"
	"net/http"
)

func GetHandler(w http.ResponseWriter, r *http.Request) {
	// read date requested from query string
	// call into data.getRidersFor5Days(date time.Time) using that date
	// return a json array of riders
}

func PutHandler(w http.ResponseWriter, r *http.Request) {
	// build up data.Rider from JSON
	// call data.AddRider(data.Rider)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	// call data.DeleteRider(data.Rider)
	// match must be exact { name, date, direction }
}
