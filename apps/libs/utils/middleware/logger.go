package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		requestID, _ := param.Keys["requestID"].(string)
		if requestID == "" {
			requestID = "unknown"
		}

		return fmt.Sprintf("[%s] [%s] %s %s %d %s \"%s\" %s \"%s\"\n",
			param.TimeStamp.Format(time.RFC3339),
			requestID,
			param.ClientIP,
			param.Method,
			param.StatusCode,
			param.Latency,
			param.Path,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	})
}
