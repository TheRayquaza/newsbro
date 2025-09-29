package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"
)

func AdminMiddleware() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		var user, exists = c.Get("user")
		if user == nil || !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Unauthorized",
				"request_id": c.GetString("requestID"),
			})
			c.Abort()
			return
		}
		u := user.(*entities.JWTClaims)
		if u.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error":      "Forbidden: Admins only",
				"request_id": c.GetString("requestID"),
			})
			c.Abort()
			return
		}
		c.Next()
	})
}
