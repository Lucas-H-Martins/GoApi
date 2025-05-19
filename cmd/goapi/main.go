package main

import (
	"log"

	"goapi/config"
	"goapi/routes"

	_ "goapi/docs" // This will be generated
)

// @title           Go API
// @version         1.0
// @description     A simple API written in Go
// @host            localhost:8080
// @BasePath        /
func main() {
	// Load configuration based on environment
	cfg, err := config.LoadConfig()

	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database
	db := config.NewPostgresDB(&cfg.Database)

	if err := db.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Setup router
	router := routes.SetupRouter()

	// Start server
	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
