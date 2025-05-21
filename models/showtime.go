package models

import (
	"gorm.io/gorm"
)

type Showtime struct {
	gorm.Model
	MovieId      uint   `json:"movie_id"`
	Time         string `json:"time"`
	Capacity     uint   `json:"capacity"`
	Reserved     uint   `json:"reserved"`
	StartTime    string `json:"start_time"`
	SeatCapacity uint   `json:"seat_capacity"`
}
