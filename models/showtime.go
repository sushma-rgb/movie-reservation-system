package models

import (
	"time"

	"gorm.io/gorm"
)

type Showtime struct {
	gorm.Model
	Id           uint      `json:"id"`
	MovieId      uint      `json:"movie_id"`
	Time         string    `json:"time"`
	Capacity     uint      `json:"capacity"`
	Reserved     uint      `json:"reserved"`
	StartTime    time.Time `json:"start_time"`
	SeatCapacity uint      `json:"seat_capacity"`
}
