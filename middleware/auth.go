package middleware

import (
	"fmt"
	"movie-reservation/config"
	"movie-reservation/models"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Just get the raw token from header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization token missing"})
			c.Abort()
			return
		}

		// Parse the token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		userIDFloat, ok := claims["id"].(float64)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
			c.Abort()
			return
		}

		role, ok := claims["role"].(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid role in token"})
			c.Abort()
			return
		}

		// Set in context
		c.Set("userID", uint(userIDFloat))
		c.Set("role", role)
		c.Next()
	}
}

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access only"})
			c.Abort()
			return
		}
		c.Next()
	}
}
func ReservationOwnerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.MustGet("userID").(uint) // from Auth middleware
		reservationIDParam := c.Param("id")  // from route param

		reservationID, err := strconv.ParseUint(reservationIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid reservation ID"})
			c.Abort()
			return
		}

		var reservation models.Reservation
		if err := config.DB.First(&reservation, reservationID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Reservation not found"})
			c.Abort()
			return
		}

		if reservation.UserId != userID {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to perform this action"})
			c.Abort()
			return
		}

		// Owner verified, continue
		c.Next()
	}
}
