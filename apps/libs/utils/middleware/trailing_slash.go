package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RemoveTrailingSlash() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		if len(path) > 1 && strings.HasSuffix(path, "/") {
			newPath := strings.TrimRight(path, "/")

			newURL := *c.Request.URL
			newURL.Path = newPath

			c.Redirect(http.StatusMovedPermanently, newURL.String())
			c.Abort()
			return
		}

		c.Next()
	})
}
