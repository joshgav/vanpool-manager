package data

import "fmt"

func getRidersFor5Days(startDate time.Time) ([]*Rider, error) {}
func addRider(r *Rider) error                                 {}
func deleteRider(r *Rider) error                              {}

type Rider struct {
	Name        string          `json:name`        // "first last"
	VanpoolName string          `json:vanpoolName` // ""
	Date        string          `json:date`        // "yyyy-mm-dd"
	Direction   TravelDirection `json:direction`   // "I" or "O"
}

type TravelDirection rune

const (
	TravelDirectionInbound  TravelDirection = "I"
	TravelDirectionOutbound TravelDirection = "O"
)
