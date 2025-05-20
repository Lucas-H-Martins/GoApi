package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a middleware for basic authentication
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "authorization header is required",
			})
			c.Abort()
			return
		}

		// For testing purposes, check if credentials match
		if authHeader != "bHVjYXNAbHVjYXMuY29tLmJyOjEyMzQ=" { // Base64 of "lucas@lucas.com.br:1234"
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "invalid credentials",
			})
			c.Abort()
			return
		}

		// Set a test user ID in the context
		c.Set("user_id", int64(1))
		c.Next()
	}
}
