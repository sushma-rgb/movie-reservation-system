package main

import (
	"movie-reservation/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Movie Reservation Backend is Running!"})
	})

	r.Run(":8080") // Start server
}
