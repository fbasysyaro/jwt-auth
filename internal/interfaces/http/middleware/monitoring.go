package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s | %s | %d | %s | %s | %s | %s | %s\n",
			param.TimeStamp.Format(time.RFC3339),
			param.Method,
			param.StatusCode,
			param.Path,
			param.ClientIP,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}

func RequestTracing() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Set request ID
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", time.Now().UnixNano())
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Log the request details
		status := c.Writer.Status()
		path := c.Request.URL.Path
		method := c.Request.Method
		userAgent := c.Request.UserAgent()

		fmt.Printf("[TRACE] %s | %s %s | %d | %v | %s | %s\n",
			requestID,
			method,
			path,
			status,
			latency,
			c.ClientIP(),
			userAgent,
		)
	}
}
