package main

import (
	"flag"
	"os"

	"goapi/config"
	"goapi/logger"
	"goapi/migrations"
)

func main() {

	// Load configuration based on environment
	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Failed to load configuration: %v", err)
		os.Exit(1)
	}
	// valid flag values
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")
	flag.Parse()

	if !*up && !*down {
		logger.Error("Please specify either -up or -down")
		os.Exit(1)
	}

	// Initialize database
	db := config.NewPostgresDB(&cfg.Database)
	// connect to the database
	if err := db.Connect(); err != nil {
		logger.Error("Failed to connect to database: %v", err)
		os.Exit(1)
	}

	// Initialize migration manager
	manager, err := migrations.NewManager(db.GetDB())
	if err != nil {
		logger.Error("Failed to create migration manager: %v", err)
		os.Exit(1)
	}

	if err := manager.LoadMigrations(); err != nil {
		logger.Error("Failed to load migrations: %v", err)
		os.Exit(1)
	}

	if *up {
		if err := manager.RunMigrationsUp(); err != nil {
			// Log the error and exit
			logger.Error("Failed to run migrations up: %v", err)
			os.Exit(1)
		}
	}

	if *down {
		if err := manager.RunMigrationsDown(); err != nil {
			// Log the error and exit
			logger.Error("Failed to run migrations down: %v", err)
			os.Exit(1)
		}
	}

	// Close the database connection
	defer db.Close()
}
