package model

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/satori/go.uuid"
)

const (
	tableName = "riders"
)

type TravelDirection string

const (
	TravelDirectionInbound  TravelDirection = "I"
	TravelDirectionOutbound TravelDirection = "O"
)

type Rider struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	Username    string          `json:"username" db:"username"`
	DisplayName string          `json:"displayName" db:"display_name"`
	Date        string          `json:"date" db:"date"`
	Direction   TravelDirection `json:"direction" db:"direction"`
}

func createTableIfNotExists() error {
	log.Printf("createTableIfNotExists()\n")
	if tableExists(tableName) {
		log.Printf("table already exists, nothing to do\n")
		return nil
	}

	db, err := database()
	if err != nil {
		log.Fatalf("failed to get open database: %s\n", err)
	}

	const q_create_table = `CREATE TABLE riders (
		id uuid,
		username varchar(100),
		display_name varchar(100),
		date date,
		direction char(1)
		);`
	_, err = db.ExecContext(context.Background(), q_create_table)
	return err
}

func tableExists(name string) bool {
	db, err := database()
	if err != nil {
		log.Fatalf("failed to get open database: %s\n", err)
	}

	rows, err := db.QueryContext(context.Background(),
		"SELECT EXISTS (SELECT 1 FROM pg_tables WHERE tablename = 'riders')")
	if err != nil {
		log.Printf("table does not yet exist: %s\n", err)
		return false
	}
	var exists bool
	rows.Next()
	_ = rows.Scan(&exists)
	log.Printf("table exists")
	return exists
}

func parseDate(date time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d", date.Year(), date.Month(), date.Day())
}

func GetRidersFor(date time.Time, direction TravelDirection) ([]Rider, error) {
	db, err := database()
	if err != nil {
		log.Fatalf("failed to get open database: %s\n", err)
	}

	_riders := []Rider{}
	err = db.SelectContext(
		context.Background(),
		&_riders,
		"SELECT * FROM riders WHERE date = $1 AND direction = $2",
		parseDate(date),
		direction,
	)
	for i, _ := range _riders {
		t, _ := time.Parse(time.RFC3339, _riders[i].Date)
		_riders[i].Date = parseDate(t)
		_riders[i].Direction = TravelDirection(_riders[i].Direction)
		log.Printf("date: %v\n", _riders[i].Date)
	}
	return _riders, err
}

func AddRider(r *Rider) error {
	db, err := database()
	if err != nil {
		log.Fatalf("failed to get open database: %s\n", err)
	}

	_, err = db.ExecContext(
		context.Background(),
		`INSERT INTO riders (id, username, display_name, date, direction)
			VALUES ($1, $2, $3, $4, $5)`,
		uuid.Must(uuid.NewV4()),
		r.Username,
		r.DisplayName,
		r.Date,
		r.Direction,
	)
	return err
}

func DeleteRider(r *Rider) error {
	db, err := database()
	if err != nil {
		log.Fatalf("failed to get open database: %s\n", err)
	}

	_, err = db.ExecContext(
		context.Background(),
		`DELETE FROM riders WHERE
			username = $1 AND date = $2 AND direction = $3`,
		r.Username,
		r.Date,
		r.Direction,
	)
	return err
}
