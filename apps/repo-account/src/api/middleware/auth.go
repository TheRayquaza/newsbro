package middleware

import (
	"net/http"
	"repo_account/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"
	"github.com/gin-gonic/gin"
)

func UserStatusMiddleware(userService *services.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userClaims, exists := c.Get("user")
		if !exists {
			c.Next()
			return
		}

		claims, ok := userClaims.(*entities.JWTClaims)
		if !ok {
			c.Next()
			return
		}

		user, err := userService.GetUserByID(claims.UserID)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		c.Next()
	}
}
