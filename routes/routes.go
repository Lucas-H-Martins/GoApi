package routes

import (
	"database/sql"
	"goapi/handlers"
	"goapi/middleware"
	"goapi/repository"
	"goapi/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// HelloWorldHandler godoc
// @Summary Returns a hello world message
// @Description Returns a simple hello world message
// @Tags hello
// @Produce json
// @Success 200 {object} map[string]string
// @Router / [get]
func HelloWorldHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello World!",
	})
}

// SetupRouter configures all the routes for the application
func SetupRouter(db *sql.DB) *gin.Engine {
	// Create a new gin router without default middleware
	router := gin.New()

	// Use our custom logger middleware
	router.Use(middleware.Logger())
	// Use recovery middleware to handle panics
	router.Use(gin.Recovery())

	// Initialize repositories
	userRepo := repository.NewPostgresUserRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.GET("/", HelloWorldHandler)

	// User routes
	router.POST("/users", userHandler.CreateUser)
	router.GET("/users", userHandler.ListUsers)
	router.GET("/users/:id", userHandler.GetUserByID)
	router.PUT("/users/:id", userHandler.UpdateUser)
	router.DELETE("/users/:id", userHandler.DeleteUser)

	return router
}
