package middleware

import (
	"net/http"
	"strings"

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

		// Check if the header has the Basic prefix
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Basic" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    http.StatusUnauthorized,
				"message": "invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Get the credentials
		credentials := parts[1]

		// For testing purposes, check if credentials match
		if credentials != "bHVjYXNAbHVjYXMuY29tLmJyOjEyMzQ=" { // Base64 of "lucas@lucas.com.br:1234"
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
