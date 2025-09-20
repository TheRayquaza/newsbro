package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func RemoveTrailingSlash() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		path := c.Request.URL.Path

		// Only redirect if path has trailing slash and is not root
		if len(path) > 1 && strings.HasSuffix(path, "/") {
			newPath := strings.TrimRight(path, "/")

			// Build new URL
			newURL := *c.Request.URL
			newURL.Path = newPath

			// Redirect to URL without trailing slash
			c.Redirect(http.StatusMovedPermanently, newURL.String())
			c.Abort()
			return
		}

		c.Next()
	})
}
