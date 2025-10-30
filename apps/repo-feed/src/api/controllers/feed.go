package controllers

import (
	"net/http"
	"strconv"
	"time"

	"repo_feed/src/api/dto"
	"repo_feed/src/converters"
	"repo_feed/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"

	"github.com/gin-gonic/gin"
)

type FeedController struct {
	feedService  *services.FeedService
	defaultModel string
}

func NewFeedController(feedService *services.FeedService, defaultModel string) *FeedController {
	return &FeedController{
		feedService:  feedService,
		defaultModel: defaultModel,
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
// @Failure 500 {object} dto.ErrorResponse
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /feed [get]
func (uc *FeedController) GetUserFeed(c *gin.Context) {
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
	if model == "" {
		model = uc.defaultModel
	}

	limit := c.Query("limit")
	if limit == "" {
		limit = "20"
	}
	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt <= 0 {
		limitInt = 20
	}
	limitInt64 := int64(limitInt)

	articles, err := uc.feedService.GetUserFeed(usr.UserID, model, limitInt64)
	if err != nil {
		switch err.(type) {
		case *dto.NotFoundError:
			c.JSON(http.StatusOK, dto.FeedResponse{
				UserID:    usr.UserID,
				ModelName: model,
				Articles:  []dto.Article{},
				UpdatedAt: time.Now(),
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			})
			return
		}
	}

	articleDTOs := make([]dto.Article, len(articles))
	for i, article := range articles {
		articleDTOs[i] = converters.ArticleModelToArticle(&article)
	}

	c.JSON(http.StatusOK, dto.FeedResponse{
		UserID:    usr.UserID,
		ModelName: model,
		Articles:  articleDTOs,
		UpdatedAt: time.Now(),
	})
}
