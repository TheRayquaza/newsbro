package controllers

import (
	"net/http"

	"repo_article/src/api/dto"
	"repo_article/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/auth/entities"

	"github.com/gin-gonic/gin"
)

type RSSController struct {
	rssService *services.RSSService
}

func NewRSSController(rssService *services.RSSService) *RSSController {
	return &RSSController{
		rssService: rssService,
	}
}

// @Summary Create a new rss
// @Description Create a new rss with the provided data
// @Tags RSS
// @Accept json
// @Produce json
// @Param rss body dto.RSSCreateRequest true "RSS data"
// @Success 201 {object} dto.RSSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss [post]
func (rc *RSSController) CreateRSS(c *gin.Context) {
	var req dto.RSSCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.JWTClaims).UserID

	rss, err := rc.rssService.CreateRSS(&req, userID)
	if err != nil {
		switch err.(type) {
		case *dto.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, rss)
}

// @Summary Get rss by ID
// @Description Get a specific rss by its ID
// @Tags RSS
// @Produce json
// @Param name path string true "RSS Name"
// @Success 200 {object} dto.RSSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss/{name} [get]
func (rc *RSSController) GetRSSByName(c *gin.Context) {
	nameParam := c.Param("name")

	rss, err := rc.rssService.GetRSSByName(nameParam)
	if err != nil {
		switch err.(type) {
		case *dto.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case nil:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rss)
}

// @Summary Get all rss
// @Description Get a list of rss with optional filtering
// @Tags RSS
// @Produce json
// @Param begin_at query string false "Filter by begin date (YYYY-MM-DD)"
// @Param end_at query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {array} dto.RSSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss [get]
func (rc *RSSController) GetRSS(c *gin.Context) {
	var filters dto.RSSFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.JWTClaims).UserID

	rss, err := rc.rssService.GetRSS(userID, &filters)
	if err != nil {
		switch err.(type) {
		case *dto.ErrBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rss)
}

// @Summary Get tree rss
// @Description Get the tree structure of rss feeds with optional filtering
// @Tags RSS
// @Produce json
// @Param begin_at query string false "Filter by begin date (YYYY-MM-DD)"
// @Param end_at query string false "Filter by end date (YYYY-MM-DD)"
// @Success 200 {array} dto.TreeRSSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss/tree [get]
func (rc *RSSController) GetTreeRSS(c *gin.Context) {
	var filters dto.RSSFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	userID := user.(*entities.JWTClaims).UserID

	rssTree, err := rc.rssService.GetTreeRSS(userID, &filters)
	if err != nil {
		switch err.(type) {
		case *dto.ErrBadRequest:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rssTree)
}

// @Summary Update rss
// @Description Update an existing rss by name
// @Tags RSS
// @Accept json
// @Produce json
// @Param name path string true "RSS Name"
// @Param rss body dto.RSSUpdateRequest true "Rss update data"
// @Success 200 {object} dto.RSSResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss/{name} [put]
func (rc *RSSController) UpdateRSS(c *gin.Context) {
	nameParam := c.Param("name")

	var req dto.RSSUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	rss, err := rc.rssService.UpdateRSS(nameParam, &req)
	if err != nil {
		switch err.(type) {
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case *dto.ErrConflict:
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, rss)
}

// @Summary Delete rss
// @Description Delete an rss by ID
// @Tags RSS
// @Param name path string true "RSS Name"
// @Success 204
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security JWT
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /rss/{name} [delete]
func (rc *RSSController) DeleteRSS(c *gin.Context) {
	nameParam := c.Param("name")

	err := rc.rssService.DeleteRSS(nameParam)
	if err != nil {
		switch err.(type) {
		case *dto.ErrNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}
