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

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"
	"github.com/TheRayquaza/newsbro/apps/libs/utils"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
)

type AuthService struct {
	Config       *config.Config
	db           *gorm.DB
	oidcVerifier *oidc.IDTokenVerifier
	oauth2Config oauth2.Config
}

func NewAuthService(cfg *config.Config, db *gorm.DB) *AuthService {
	service := &AuthService{
		Config: cfg,
		db:     db,
	}

	if cfg.OIDCIssuerURL != "" {
		service.initOIDC()
	}

	return service
}

func (s *AuthService) initOIDC() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, s.Config.OIDCIssuerURL)
	if err != nil {
		utils.SugarLog.Errorf("Failed to initialize OIDC provider: %v\n", err)
		return
	}

	s.oidcVerifier = provider.Verifier(&oidc.Config{
		ClientID: s.Config.OIDCClientID,
	})

	s.oauth2Config = oauth2.Config{
		ClientID:     s.Config.OIDCClientID,
		ClientSecret: s.Config.OIDCClientSecret,
		RedirectURL:  s.Config.OIDCRedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}
}

func (s *AuthService) Register(req *dto.RegisterRequest) (*dto.LoginResponse, error) {
	var existingUser models.User
	if err := s.db.Where("email = ? OR username = ?", req.Email, req.Username).First(&existingUser).Error; err == nil {
		utils.SugarLog.Warn("Attempt to register with existing email or username:", req.Email, req.Username)
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

	utils.SugarLog.Debug("Registering new user:", req.Email, req.Username)
	if err := user.SetPassword(req.Password); err != nil {
		return nil, err
	}

	utils.SugarLog.Info("Creating user record in database for:", req.Email)
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return s.generateTokens(&user)
}

func (s *AuthService) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
	var user models.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		utils.SugarLog.Warn("Login attempt with non-existent email:", req.Email)
		return nil, errors.New("invalid credentials")
	}

	if !user.CheckPassword(req.Password) {
		utils.SugarLog.Warn("Invalid password attempt for email:", req.Email)
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		utils.SugarLog.Warn("Attempt to login to inactive account:", user.Email)
		return nil, errors.New("account is inactive")
	}

	return s.generateTokens(&user)
}

func (s *AuthService) RefreshToken(user *models.User) (*dto.LoginResponse, error) {
	var token *models.RefreshToken
	if err := s.db.Preload("User").Where("user_id = ? AND expires_at > ?", user.ID, time.Now()).First(&token).Error; err != nil {
		utils.SugarLog.Warn("Invalid or expired refresh token for user ID:", user.ID)
		return nil, errors.New("invalid refresh token")
	}

	utils.SugarLog.Info("Refreshing tokens for user ID:", user.ID)
	s.db.Delete(&token)

	return s.generateTokens(&token.User)

}

func (s *AuthService) generateTokens(user *models.User) (*dto.LoginResponse, error) {
	accessToken, err := s.GenerateAccessToken(user)
	if err != nil {
		utils.SugarLog.Errorf("Error generating access token for user ID: %d, Error: %v", user.ID, err)
		return nil, err
	}

	refreshTokenStr, err := s.generateRefreshToken(user)
	if err != nil {
		utils.SugarLog.Errorf("Error generating refresh token for user ID: %d, Error: %v", user.ID, err)
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
	claims := entities.JWTClaims{
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
	return token.SignedString([]byte(s.Config.JWTSecret))
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
		utils.SugarLog.Errorf("Error saving refresh token for user ID: %d, Error: %v", user.ID, err)
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
		utils.SugarLog.Warn("OIDC not configured")
		return nil, errors.New("OIDC not configured")
	}

	ctx := context.Background()

	utils.SugarLog.Debug("Exchanging code for token with OIDC provider")
	oauth2Token, err := s.oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, err
	}

	utils.SugarLog.Debug("Verifying ID token")
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}

	utils.SugarLog.Debug("Parsing ID token claims")
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

	utils.SugarLog.Info("OIDC login for email:", claims.Email)

	var user models.User
	err = s.db.Where("email = ?", claims.Email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.SugarLog.Info("Creating new user from OIDC claims for email:", claims.Email)
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
				utils.SugarLog.Errorf("Error creating user from OIDC claims: %v", err)
				return nil, err
			}
		} else {
			utils.SugarLog.Errorf("Error fetching user by email: %v", err)
			return nil, err
		}
	}

	return s.generateTokens(&user)
}

func (s *AuthService) GetPostOAuthRedirectURL() string {
	return s.Config.OIDCRedirectFrontendURL
}
