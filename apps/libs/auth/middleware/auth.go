package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if cookie, err := c.Request.Cookie("auth_token"); err == nil && cookie.Value != "" {
			token := cookie.Value
			if authHeader != "" && authHeader != "Bearer "+token {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error":      "Conflicting authentication methods",
					"request_id": c.GetString("requestID"),
				})
				c.Abort()
				return
			}
			authHeader = "Bearer " + token
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

		user, err := validateToken(tokenParts[1], jwtSecret)
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

func validateToken(tokenStr string, jwtSecret string) (*entities.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &entities.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*entities.JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
