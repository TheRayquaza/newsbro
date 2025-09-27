package middleware

import (
	"net/http"
	"strings"

	"repo_article/src/domain/services"
	"repo_article/src/config"
	"repo_article/src/domain/entities"

	"github.com/gin-gonic/gin"
)

func AdminMiddleware(authService *services.AuthService) gin.HandlerFunc {
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
		u := user.(*entities.User)
		isAdmin := false
		for _, admin := range config.Load().Admins {
			if strings.EqualFold(u.Email, admin) {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
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
