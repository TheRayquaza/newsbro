package controllers

import (
	"net/http"

	"repo_feed/src/api/dto"
	"repo_feed/src/converters"
	"repo_feed/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
	feedService *services.FeedService
}

func NewFeedController(feedService *services.FeedService) *FeedController {
	return &FeedController{
		feedService: feedService,
	}
}

// FeedController - Add model query parameter
// @Summary Get user feed
// @Description Get the feed of the authenticated user
// @Tags Feed
// @Produce json
// @Param model query string false "Model name (optional, uses default if not specified)"
// @Success 200 {object} dto.FeedResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /feed [get]
func (uc *FeedController) GetFeed(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
			Code:    http.StatusUnauthorized,
			Message: "User not found in context",
		})
		return
	}

	usr, ok := user.(*entities.JWTClaims)
	if !ok {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Invalid user type in context",
		})
		return
	}

	model := c.Query("model")

	feed, err := uc.feedService.GetFeed(usr.UserID, model)
	if err != nil {
		if err.Error() == "feed not found" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Code:    http.StatusNotFound,
				Message: "Feed not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to retrieve feed",
		})
		return
	}

	c.JSON(http.StatusOK, converters.FeedModelToFeedResponse(feed))
}
