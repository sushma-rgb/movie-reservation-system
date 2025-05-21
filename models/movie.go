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

type Input struct {
	Id          int
	Title       string
	Description string
	Genre       string
	Poster      string
	Showtimes   []Showtime
}
