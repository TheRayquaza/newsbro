package controllers

import (
	"net/http"
	"strconv"

	"repo_account/src/api/dto"
	"repo_account/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

// @Summary Get user profile
// @Description Get the profile of the authenticated user
// @Tags User
// @Produce json
// @Success 200 {object} dto.UserResponse
// @Failure 401
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /users/profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	usr, ok := user.(*entities.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	c.JSON(http.StatusOK, dto.UserResponse{
		ID:        usr.UserID,
		Email:     usr.Email,
		Username:  usr.Username,
		FirstName: usr.FirstName,
		LastName:  usr.LastName,
		Role:      usr.Role,
		IsActive:  true,
	})
}

// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param updateRequest body map[string]interface{} true "Update Request"
// @Success 200 {object} dto.UserResponse
// @Failure 400
// @Failure 401
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /users/profile [put]
func (uc *UserController) UpdateProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Remove sensitive fields from updates
	delete(updates, "id")
	delete(updates, "password")
	delete(updates, "username")

	usr, ok := user.(*entities.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	updatedUser, err := uc.userService.UpdateUser(usr.UserID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := dto.UserResponse{
		ID:        updatedUser.ID,
		Email:     updatedUser.Email,
		Username:  updatedUser.Username,
		FirstName: updatedUser.FirstName,
		LastName:  updatedUser.LastName,
		Role:      updatedUser.Role,
		IsActive:  updatedUser.IsActive,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Get all users
// @Description Get a list of all users (admin only)
// @Tags User
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} dto.UserResponse
// @Failure 403
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /users [get]
func (uc *UserController) GetUsers(c *gin.Context) {
	limit := 10
	offset := 0

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil {
			limit = parsed
		}
	}

	if o := c.Query("offset"); o != "" {
		if parsed, err := strconv.Atoi(o); err == nil {
			offset = parsed
		}
	}

	users, err := uc.userService.GetAllUsers(limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	response := make([]dto.UserResponse, len(users))
	for i, usr := range users {
		response[i] = dto.UserResponse{
			ID:        usr.ID,
			Email:     usr.Email,
			Username:  usr.Username,
			FirstName: usr.FirstName,
			LastName:  usr.LastName,
			Role:      usr.Role,
			IsActive:  usr.IsActive,
		}
	}

	c.JSON(http.StatusOK, response)
}
