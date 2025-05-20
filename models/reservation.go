package models

import "gorm.io/gorm"

type Reservation struct {
	gorm.Model
	UserId     uint   `json:"user_id"`
	ShowtimeId uint   `json:"showtime_id"`
	Seats      uint   `json:"seats"`
	Status     string `json:"status"`
}
