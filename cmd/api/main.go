package main

import (
	"go-article/internal/config"
	"go-article/internal/routes"
	"os"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	config.ConnectDatabase()

	// Setup Router
	r := routes.SetupRoutes(config.DB)

	// Run Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
