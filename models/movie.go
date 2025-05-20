package models

import "gorm.io/gorm"

type Movie struct {
	gorm.Model
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Genre       string     `json:"genre"`
	Poster      string     `json:"poster"`
	Showtimes   []Showtime `json:"showtimes"`
}
