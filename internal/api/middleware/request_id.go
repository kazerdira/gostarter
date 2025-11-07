package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestID generates a unique ID for each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")
		
		// Generate new ID if not provided
		if requestID == "" {
			requestID = uuid.New().String()
		}

		// Store in context
		c.Set("request_id", requestID)

		// Add to response header
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
