package routes

import (
	"movie-reservation/controllers"
	"movie-reservation/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	// Public routes
	r.POST("/signup", controllers.Signup)
	r.POST("/login", controllers.Login)

	// Authenticated user routes
	userGroup := r.Group("/user")
	userGroup.Use(middleware.AuthMiddleware())
	{
		userGroup.GET("/me", controllers.GetProfile)
		userGroup.POST("/reserve", controllers.ReserveSeats)
		userGroup.GET("/showtimes/:id/available", controllers.GetAvailableSeats)
		userGroup.DELETE("/reservation/:id/cancel", controllers.CancelReservation)
		userGroup.DELETE("/reservation/:id/cancel", middleware.ReservationOwnerMiddleware(), controllers.CancelReservation)
	}

	adminGroup := r.Group("/admin")
	adminGroup.Use(middleware.AuthMiddleware(), middleware.AdminMiddleware())
	{
		// Movie management
		adminGroup.POST("/movies", controllers.AddMovie)
		adminGroup.PUT("/movies/:id", controllers.UpdateMovie)
		adminGroup.DELETE("/movies/:id", controllers.DeleteMovie)

		// Showtime management
		adminGroup.POST("/showtimes", controllers.AddShowtime)
		adminGroup.PUT("/showtimes/:id", controllers.UpdateShowtime)
		adminGroup.DELETE("/showtimes/:id", controllers.DeleteShowtime)
	}
}
