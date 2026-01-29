package config

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// RequestLoggerMiddleware logs the HTTP request details
func RequestLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// Process the request
		c.Next()

		// Log after processing, but before the connection closes technically
		// Ideally, you'd use a structured logger like zap or zerolog here
		latency := time.Since(startTime)
		status := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.URL.Path

		log.Printf("[HTTP] %d | %12v | %s %s", status, latency, method, path)
	}
}
