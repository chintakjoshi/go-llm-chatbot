package middleware

import (
	"chintak-chatbot/utils"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get client IP
		clientIP := c.ClientIP()

		// Log the request
		utils.LogAPIRequest(
			c.Request.Method,
			c.Request.URL.Path,
			clientIP,
			c.Writer.Status(),
			duration,
		)

		// Log slow requests
		if duration > time.Second {
			utils.Warn("Slow Request - %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}

		// Log errors
		if c.Writer.Status() >= 400 {
			utils.Error("HTTP Error - %s %s - Status: %d - Client: %s",
				c.Request.Method, c.Request.URL.Path, c.Writer.Status(), clientIP)
		}
	}
}
