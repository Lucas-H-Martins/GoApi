package middleware

import (
	"time"

	"goapi/logger"

	"github.com/gin-gonic/gin"
)

// Logger returns a gin middleware that logs requests using our custom logger
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get method
		method := c.Request.Method

		// Get error message if any
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		// Log the request
		if statusCode >= 500 {
			logger.Error("[GIN] %s | %3d | %13v | %15s | %s %s %s",
				method,
				statusCode,
				latency,
				clientIP,
				path,
				raw,
				errorMessage,
			)
		} else if statusCode >= 400 {
			logger.Warn("[GIN] %s | %3d | %13v | %15s | %s %s %s",
				method,
				statusCode,
				latency,
				clientIP,
				path,
				raw,
				errorMessage,
			)
		} else {
			logger.Info("[GIN] %s | %3d | %13v | %15s | %s %s",
				method,
				statusCode,
				latency,
				clientIP,
				path,
				raw,
			)
		}
	}
}
