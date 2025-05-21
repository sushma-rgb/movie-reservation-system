package main

import (
	"movie-reservation/config"
	"movie-reservation/routes"
)

func main() {
	config.ConnectDB()

	r := routes.SetupRoutes()

	r.Run(":8080")

}
