package routes

import (
	"goapi/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRouter configures all the routes for the application
func SetupRouter() *gin.Engine {
	// Create a new gin router without default middleware
	router := gin.New()

	// Use our custom logger middleware
	router.Use(middleware.Logger())
	// Use recovery middleware to handle panics
	router.Use(gin.Recovery())

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Hello World route
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// router.GET("/users", userHandler.GetUsers)
	// router.GET("/users/:id", userHandler.GetUserByID)

	return router
}
