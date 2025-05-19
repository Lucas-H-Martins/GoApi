package routes

import (
	"goapi/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter configures all the routes for the application
func SetupRouter() *gin.Engine {
	// Create a new gin router without default middleware
	router := gin.New()

	// Use our custom logger middleware
	router.Use(middleware.Logger())
	// Use recovery middleware to handle panics
	router.Use(gin.Recovery())

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
