package controllers

import (
	"movie-reservation/config"
	"movie-reservation/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Admin can add a movie
func AddMovie(c *gin.Context) {
	var movie models.Movie
	if err := c.Bind(&movie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Create(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add movie"})
		return
	}

	c.JSON(http.StatusCreated, movie)
}

// Admin can update a movie by ID
func UpdateMovie(c *gin.Context) {
	id := c.Param("id")
	var movie models.Movie

	if err := config.DB.First(&movie, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		return
	}

	var input models.Movie
	if err := c.Bind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	movie.Title = input.Title
	movie.Description = input.Description
	movie.Genre = input.Genre
	movie.Poster = input.Poster

	if err := config.DB.Save(&movie).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// Admin can delete a movie by ID
func DeleteMovie(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Movie{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete movie"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Movie deleted successfully"})
}
