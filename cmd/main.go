package main

import (
	"log"

	"mock-ses-api/internal/config"
	"mock-ses-api/internal/routes"

	"github.com/joho/godotenv"
)

func main() {

	// Load the environment variables
	godotenv.Load("../.env")

	// Initialize the configuration
	cfg := config.New()

	// Initialize the router
	router := routes.SetupRouter(cfg)

	log.Printf("Server starting on %s", cfg.ServerAddress)

	// Start the server
	if err := router.Run(cfg.Domain + ":" + cfg.ServerAddress); err != nil {
		log.Fatal(err)
	}
}
