package routes

import (
	"database/sql"
	"goapi/middleware"
	"goapi/routes/user_routes"

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

	// Swagger documentation - placed before auth middleware to be publicly accessible
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Use our custom authorization middleware
	router.Use(gin.HandlerFunc(middleware.AuthMiddleware()))

	// Setup user routes
	user_routes.SetupUserRoutes(router, db)

	return router
}
