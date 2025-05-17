package routes

import (
	"github.com/gorilla/mux"
	"goapi/controllers"
	"goapi/middleware"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Public routes (no authentication required)
	router.HandleFunc("/hello", controllers.HelloWorld).Methods("GET")

	// Protected routes (require authentication)
	protected := router.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware)
	
	// Add protected endpoints here
	protected.HandleFunc("/protected/hello", controllers.HelloWorld).Methods("GET")

	return router
} 