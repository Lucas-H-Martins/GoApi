package main

import (
	"log"

	"goapi/config"
	// "goapi/routes"
)

func main() {
	// Load configuration based on environment
	_, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// // Initialize database
	// db := config.NewPostgresDB(&cfg.Database)
	// if err := db.Connect(); err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }
	// defer db.Close()

	// // Initialize router
	// router := routes.NewRouter(db.GetDB())

	// // Configure server
	// addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	// server := &http.Server{
	// 	Addr:    addr,
	// 	Handler: router,
	// }

	// // Log startup information
	// log.Printf("Starting server in %s mode", cfg.Environment)
	// log.Printf("Server is running on %s", addr)

	// // Start server
	// if err := server.ListenAndServe(); err != nil {
	// 	log.Fatalf("Server failed to start: %v", err)
	// }
}
