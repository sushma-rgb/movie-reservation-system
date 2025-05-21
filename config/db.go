package config

import (
	"log"
	"movie-reservation/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=localhost user=postgres password=postgres dbname=movie_db port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database: ", err)
	}

	log.Println("Connected to the database ✅")

	database.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Showtime{},
		&models.Reservation{},
	)

	DB = database
}
