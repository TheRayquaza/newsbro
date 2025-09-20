package controllers

import (
	"net/http"
	"strconv"

	"repo_account/src/data/models"
	"repo_account/src/domain/services"

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
// @Success 200 {object} models.User
// @Failure 401
// @Router /user/profile [get]
func (uc *UserController) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	c.JSON(http.StatusOK, user.(*models.User))
}

// @Summary Update user profile
// @Description Update the profile of the authenticated user
// @Tags User
// @Accept json
// @Produce json
// @Param updateRequest body map[string]interface{} true "Update Request"
// @Success 200 {object} models.User
// @Failure 400
// @Failure 401
// @Router /user/profile [put]
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
	delete(updates, "oidc_subject")

	updatedUser, err := uc.userService.UpdateUser(user.(*models.User).ID, updates)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// @Summary Get all users
// @Description Get a list of all users (admin only)
// @Tags User
// @Produce json
// @Param limit query int false "Limit" default(10)
// @Param offset query int false "Offset" default(0)
// @Success 200 {array} models.User
// @Failure 403
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

	c.JSON(http.StatusOK, gin.H{"users": users})
}
