package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"

	"github.com/gin-gonic/gin"
	"libs/auth/services"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		session := sessions.Default(c)
		cookie, exists := session.Get("auth_token").(string)
		if exists && cookie != "" {
			if authHeader != "" && authHeader != "Bearer "+cookie {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":      "Conflicting authentication methods",
					"request_id": c.GetString("requestID"),
				})
				c.Abort()
				return
			}
			authHeader = "Bearer " + cookie
		}
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Authorization header required",
				"request_id": c.GetString("requestID"),
			})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Invalid authorization header format",
				"request_id": c.GetString("requestID"),
			})
			c.Abort()
			return
		}

		user, err := authService.ValidateToken(tokenParts[1])
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":      "Invalid token",
				"request_id": c.GetString("requestID"),
			})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	})
}
