package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"repo_account/src/api/dto"
	"repo_account/src/config"
	"repo_account/src/data/models"

	"github.com/coreos/go-oidc/v3/oidc"
	libJwt "github.com/TheRayquaza/newsbro/apps/libs/auth/services"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService struct {
	config       *config.Config
	db           *gorm.DB
	oidcVerifier *oidc.IDTokenVerifier
	oauth2Config oauth2.Config
}

type JWTClaims struct {
	UserID    uint   `json:"user_id"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Role      string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(cfg *config.Config, db *gorm.DB) *AuthService {
	service := &AuthService{
		config: cfg,
		db:     db,
	}

	if cfg.OIDCIssuerURL != "" {
		service.initOIDC()
	}

	return service
}

func (s *AuthService) initOIDC() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, s.config.OIDCIssuerURL)
	if err != nil {
		fmt.Printf("Failed to initialize OIDC provider: %v\n", err)
		return
	}

	s.oidcVerifier = provider.Verifier(&oidc.Config{
		ClientID: s.config.OIDCClientID,
	})

	s.oauth2Config = oauth2.Config{
		ClientID:     s.config.OIDCClientID,
		ClientSecret: s.config.OIDCClientSecret,
		RedirectURL:  s.config.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Check if user already exists
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("user already exists")
	}

	user := models.User{
		Email:     req.Email,
		Username:  req.Username,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
		Role:      "user",
	}

	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return s.generateTokens(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	return s.generateTokens(&user)
}

func (s *AuthService) RefreshToken(user *models.User) (*dto.LoginResponse, error) {
	var token *models.RefreshToken
	if err := s.db.Preload("User").Where("user_id = ? AND expires_at > ?", user.ID, time.Now()).First(&token).Error; err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Delete old refresh token
	s.db.Delete(&token)

	return s.generateTokens(&token.User)

}

func (s *AuthService) generateTokens(user *models.User) (*dto.LoginResponse, error) {
	// Generate access token
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshTokenStr, err := s.generateRefreshToken(user)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshTokenStr,
		ExpiresIn:    3600, // 1 hour
		User: dto.UserResponse{
			ID:        user.ID,
			Email:     user.Email,
			Username:  user.Username,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Role:      user.Role,
			IsActive:  user.IsActive,
		},
	}, nil
}

func (s *AuthService) GenerateAccessToken(user *models.User) (string, error) {
	claims := JWTClaims{
		UserID:    user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Username:  user.Username,
		Role:      user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "repo-account",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWTSecret))
}

func (s *AuthService) generateRefreshToken(user *models.User) (string, error) {
	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		return "", err
	}
	tokenStr := hex.EncodeToString(tokenBytes)

	refreshToken := models.RefreshToken{
		Token:     tokenStr,
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(24 * 7 * time.Hour), // 7 days
	}

	if err := s.db.Create(&refreshToken).Error; err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (s *AuthService) GetOAuthURL(state string) string {
	if s.oauth2Config.ClientID == "" {
		return ""
	}
	return s.oauth2Config.AuthCodeURL(state)
}

func (s *AuthService) HandleOAuthCallback(code string) (*dto.LoginResponse, error) {
	if s.oidcVerifier == nil {
		return nil, errors.New("OIDC not configured")
	}

	ctx := context.Background()

	// Exchange code for token
	oauth2Token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	// Verify ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}

	idToken, err := s.oidcVerifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, err
	}

	var claims struct {
		Email      string `json:"email"`
		Name       string `json:"name"`
		Subject    string `json:"sub"`
		GivenName  string `json:"given_name"`
		FamilyName string `json:"family_name"`
		Role       string `json:"role"`
		Iss        string `json:"iss"`
	}

	if err := idToken.Claims(&claims); err != nil {
		return nil, err
	}

	// Find or create user
	var user models.User
	err = s.db.Where("email = ?", claims.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Create new user
			user = models.User{
				Email:     claims.Email,
				Username:  claims.Subject,
				FirstName: claims.GivenName,
				LastName:  claims.FamilyName,
				Role:      claims.Role,
				Issuer:    &claims.Iss,
				IsActive:  true,
			}
			if err := s.db.Create(&user).Error; err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return s.generateTokens(&user)
}
