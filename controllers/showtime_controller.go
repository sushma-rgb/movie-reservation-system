package controllers

import (
	"movie-reservation/config"
	"movie-reservation/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin can add showtime for a movie
func AddShowtime(c *gin.Context) {
	var showtime models.Showtime
	if err := c.ShouldBindJSON(&showtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Optional: Validate movie exists
	var movie models.Movie
	if err := config.DB.First(&movie, showtime.MovieId).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Movie does not exist"})
		return
	}

	if err := config.DB.Create(&showtime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add showtime"})
		return
	}

	c.JSON(http.StatusCreated, showtime)
}

// Admin can update showtime
func UpdateShowtime(c *gin.Context) {
	id := c.Param("id")
	var showtime models.Showtime

	if err := config.DB.First(&showtime, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	var input models.Showtime
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	showtime.StartTime = input.StartTime
	showtime.MovieId = input.MovieId
	showtime.SeatCapacity = input.SeatCapacity

	if err := config.DB.Save(&showtime).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update showtime"})
		return
	}

	c.JSON(http.StatusOK, showtime)
}

// Admin can delete showtime by ID
func DeleteShowtime(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Showtime{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete showtime"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Showtime deleted successfully"})
}
