package data

import "testing"

func TestConnection(t *testing.T) {
	db, err := database()
	if (db == nil) && (err == nil) {
		t.Error("both db and err should not be nil")
	}
	// to test cached singleton
	db, err = database()
}
