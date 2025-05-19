package main

import (
	"log"

	"goapi/config"
	// "goapi/routes"
)

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

}
