package data

import (
	"github.com/satori/go.uuid"
	"log"
	"testing"
	"time"
)

func TestTableWorks(t *testing.T) {
	const (
		date      = "2018-01-01"
		direction = TravelDirectionInbound
	)

	var err error
	err = createTableIfNotExists()
	if err != nil {
		t.Errorf("failed to create table: %s\n", err)
	}

	rider1 := &Rider{
		ID:          uuid.Must(uuid.NewV4()),
		Username:    "testrider1@test.com",
		DisplayName: "Test Rider 1",
		Date:        date,
		Direction:   direction,
	}
	err = addRider(rider1)
	if err != nil {
		t.Errorf("failed to add rider: %s\n", err)
	}

	rider2 := &Rider{
		ID:          uuid.Must(uuid.NewV4()),
		Username:    "testrider2@test.com",
		DisplayName: "Test Rider 2",
		Date:        date,
		Direction:   direction,
	}
	err = addRider(rider2)
	if err != nil {
		t.Errorf("failed to add rider: %s\n", err)
	}

	d, _ := time.Parse("2006-01-02", date)

	riders, err := getRidersFor(d, direction)
	if err != nil {
		t.Errorf("failed to get riders: %s\n", err)
	} else {
		log.Printf("successfully got riders %v\n", riders)
	}

	if riders[0].Username != rider1.Username {
		t.Errorf("username doesn't match: %v != %v", riders[0].Username, rider1.Username)
	}

	if riders[0].DisplayName != rider1.DisplayName {
		t.Errorf("displayName doesn't match: %v != %v", riders[0].DisplayName, rider1.DisplayName)
	}

	if riders[0].Date != rider1.Date {
		t.Errorf("date doesn't match: %v != %v", riders[0].Date, rider1.Date)
	}

	if riders[0].Direction != rider1.Direction {
		t.Errorf("direction doesn't match: %v != %v", riders[0].Direction, rider1.Direction)
	}
}
