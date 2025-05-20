package user_routes

import (
	"database/sql"

	"goapi/handlers"
	"goapi/repository"
	"goapi/services"

	"github.com/gin-gonic/gin"
)

// SetupUserRoutes configures all user-related routes
func SetupUserRoutes(router *gin.Engine, db *sql.DB) {
	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// User routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.ListUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)
}
