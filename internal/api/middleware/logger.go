package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Logger returns a middleware that logs HTTP requests
func Logger(logger zerolog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log after request is processed
		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		errorMessage := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		// Build log entry
		logEvent := logger.Info()

		if statusCode >= 500 {
			logEvent = logger.Error()
		} else if statusCode >= 400 {
			logEvent = logger.Warn()
		}

		logEvent.
			Str("method", method).
			Str("path", path).
			Int("status", statusCode).
			Dur("latency", latency).
			Str("client_ip", clientIP).
			Str("request_id", c.GetString("request_id"))

		if errorMessage != "" {
			logEvent.Str("error", errorMessage)
		}

		logEvent.Msg("HTTP request")
	}
}

// Recovery returns a middleware that recovers from panics
func Recovery(logger zerolog.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logger.Error().
			Interface("error", recovered).
			Str("request_id", c.GetString("request_id")).
			Msg("Panic recovered")

		c.JSON(500, gin.H{
			"error": "internal server error",
		})
	})
}
