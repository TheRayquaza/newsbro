package services

import (
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"

	"repo_article/src/domain/entities"
)

type AuthService struct{}

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	jwt.RegisteredClaims
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) ValidateToken(tokenStr string) (*entities.User, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return nil, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return &entities.User{
		ID:        claims.UserID,
		Email:     claims.Email,
		FirstName: claims.FirstName,
		LastName:  claims.LastName,
		Username:  claims.Username,
	}, nil
}
