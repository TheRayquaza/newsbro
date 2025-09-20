package middleware

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/gin-gonic/gin"
)

func RequestID() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		// Check if request ID already exists in header
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			// Generate new request ID
			bytes := make([]byte, 16)
			if _, err := rand.Read(bytes); err == nil {
				requestID = hex.EncodeToString(bytes)
			} else {
				requestID = "unknown"
			}
		}

		// Set request ID in context and response header
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)

		c.Next()
	})
}
