package controllers

import (
	"movie-reservation/config"
	"movie-reservation/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ReserveRequest struct {
	ShowtimeId uint `json:"showtime_id"`
	Seats      uint `json:"seats"`
}

func ReserveSeats(c *gin.Context) {
	var req ReserveRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.MustGet("userID").(uint)

	var showtime models.Showtime
	if err := config.DB.First(&showtime, req.ShowtimeId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	// Check availability
	if showtime.Reserved+req.Seats > showtime.Capacity {
		c.JSON(http.StatusConflict, gin.H{"error": "Not enough available seats"})
		return
	}

	// Use transaction for atomic update and reservation creation
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Update reserved count
		if err := tx.Model(&showtime).
			UpdateColumn("reserved", gorm.Expr("reserved + ?", req.Seats)).Error; err != nil {
			return err
		}

		// Create reservation
		reservation := models.Reservation{
			UserId:     userID,
			ShowtimeId: req.ShowtimeId,
			Seats:      req.Seats,
			Status:     "active",
		}
		if err := tx.Create(&reservation).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Seats reserved successfully"})
}
func GetAvailableSeats(c *gin.Context) {
	showtimeID := c.Param("id")

	var showtime models.Showtime
	if err := config.DB.First(&showtime, showtimeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		return
	}

	availableSeats := showtime.Capacity - showtime.Reserved

	c.JSON(http.StatusOK, gin.H{
		"showtime_id":     showtime.ID,
		"available_seats": availableSeats,
	})
}
func CancelReservation(c *gin.Context) {
	reservationID := c.Param("id")
	userID := c.MustGet("userID").(uint)

	var reservation models.Reservation
	if err := config.DB.First(&reservation, reservationID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
		return
	}

	// Only reservation owner can cancel
	if reservation.UserId != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized"})
		return
	}

	if reservation.Status != "active" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Reservation already cancelled"})
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {
		// Mark reservation as cancelled
		if err := tx.Model(&reservation).Update("status", "cancelled").Error; err != nil {
			return err
		}

		// Decrease reserved seats in showtime
		if err := tx.Model(&models.Showtime{}).Where("id = ?", reservation.ShowtimeId).
			UpdateColumn("reserved", gorm.Expr("reserved - ?", reservation.Seats)).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel reservation"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reservation cancelled successfully"})
}
