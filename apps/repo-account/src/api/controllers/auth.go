package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"repo_account/src/api/dto"
	"repo_account/src/data/models"
	"repo_account/src/domain/services"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
	}
}

// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param registerRequest body dto.RegisterRequest true "Register Request"
// @Success 201
// @Failure 400 {object} dto.ErrorResponse
// @Router /auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := ac.authService.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param loginRequest body dto.LoginRequest true "Login Request"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	response, err := ac.authService.Login(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	c.SetCookie("auth_token", response.AccessToken, 3600, "/", ac.authService.Config.CookieDomain, true, true)
	c.SetCookie("refresh_token", response.RefreshToken, 86400, "/", ac.authService.Config.CookieDomain, true, true)

	c.JSON(http.StatusOK, response)
}

// @Summary Refresh JWT token
// @Description Refresh JWT token using refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /auth/refresh [post]
func (ac *AuthController) RefreshToken(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	usr, ok := user.(*models.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	response, err := ac.authService.RefreshToken(usr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary OAuth Login
// @Description Redirect to OAuth provider for login
// @Tags Auth
// @Produce json
// @Success 302
// @Failure 503 {object} dto.ErrorResponse
// @Router /auth/oauth/login [get]
func (ac *AuthController) OAuthLogin(c *gin.Context) {
	stateBytes := make([]byte, 16)
	rand.Read(stateBytes)
	state := hex.EncodeToString(stateBytes)

	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

	authURL := ac.authService.GetOAuthURL(state)
	if authURL == "" {
		c.JSON(http.StatusServiceUnavailable, dto.ErrorResponse{
			Code:    http.StatusServiceUnavailable,
			Message: "OAuth provider not configured",
		})
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// @Summary OAuth Callback
// @Description Handle OAuth provider callback
// @Tags Auth
// @Produce json
// @Param code query string true "Authorization Code"
// @Param state query string true "State Parameter"
// @Success 302
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /auth/oauth/callback [get]
func (ac *AuthController) OAuthCallback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")

	storedState, err := c.Cookie("oauth_state")
	if err != nil || storedState != state {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid state parameter",
		})
		return
	}

	c.SetCookie("oauth_state", "", -1, "/", "", false, true)

	response, err := ac.authService.HandleOAuthCallback(code)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		})
		return
	}

	c.SetCookie("auth_token", response.AccessToken, 3600, "/", ac.authService.Config.CookieDomain, true, true)
	c.SetCookie("refresh_token", response.RefreshToken, 86400, "/", ac.authService.Config.CookieDomain, true, true)

	c.Redirect(http.StatusFound, ac.authService.GetPostOAuthRedirectURL())
}
