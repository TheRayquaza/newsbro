package controllers

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"repo_article/src/api/dto"
	"repo_article/src/domain/entities"
	"repo_article/src/domain/services"

	"github.com/gin-gonic/gin"
)

type FeedbackController struct {
	feedbackService *services.FeedbackService
}

func NewFeedbackController(feedbackService *services.FeedbackService) *FeedbackController {
	return &FeedbackController{
		feedbackService: feedbackService,
	}
}

// GetArticleFeedback godoc
// @Summary Get feedback for a specific article
// @Description Get feedback statistics and user's feedback for a specific article
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Security BearerAuth
// @Success 200 {object} dto.FeedbackStatsResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /articles/{id}/feedback [get]
func (fc *FeedbackController) GetArticleFeedback(c *gin.Context) {
	newsID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid article ID")
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := user.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	userID := usr.ID

	stats, userFeedback, err := fc.feedbackService.GetArticleFeedback(uint(newsID), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	response := gin.H{
		"stats":         stats,
		"user_feedback": userFeedback,
	}

	c.JSON(http.StatusOK, response)
}

// CreateFeedback godoc
// @Summary Create feedback for an article
// @Description Create or update user feedback for a specific article
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param feedback body dto.FeedbackRequest true "Feedback data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /articles/{id}/feedback [post]
func (fc *FeedbackController) CreateFeedback(c *gin.Context) {
	newsID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid article ID")
		return
	}

	var req dto.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := user.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	userID := usr.ID

	feedback, err := fc.feedbackService.CreateOrUpdateFeedback(userID, uint(newsID), req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusCreated, feedback)
}

// UpdateFeedback godoc
// @Summary Update feedback for an article
// @Description Update existing user feedback for a specific article
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Param feedback body dto.FeedbackRequest true "Updated feedback data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /articles/{id}/feedback [put]
func (fc *FeedbackController) UpdateFeedback(c *gin.Context) {
	newsID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid article ID")
		return
	}

	var req dto.FeedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := user.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	userID := usr.ID

	_, err = fc.feedbackService.UpdateFeedback(userID, uint(newsID), req.Value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// DeleteFeedback godoc
// @Summary Delete feedback for an article
// @Description Delete user feedback for a specific article
// @Tags feedback
// @Accept json
// @Produce json
// @Param id path int true "Article ID"
// @Security BearerAuth
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /articles/{id}/feedback [delete]
func (fc *FeedbackController) DeleteFeedback(c *gin.Context) {
	newsID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Invalid article ID")
		return
	}

	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}
	usr, ok := user.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	userID := usr.ID

	err = fc.feedbackService.DeleteFeedback(userID, uint(newsID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.Status(http.StatusNoContent)
}

// GetUserFeedback godoc
// @Summary Get user's feedback history
// @Description Get all feedback given by the authenticated user
// @Tags feedback
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {array} models.Feedback
// @Router /feedback/my-feedback [get]
func (fc *FeedbackController) GetUserFeedback(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found in context"})
		return
	}

	usr, ok := user.(*entities.User)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type in context"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	feedback, _, err := fc.feedbackService.GetUserFeedback(usr.ID, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	/*
		response := gin.H{
			"data":       feedback,
			"pagination": utils.GetPaginationInfo(page, limit, int(total)),
		}
	*/

	c.JSON(http.StatusOK, feedback)
}

// ExportFeedbackCSV godoc
// @Summary Export all feedback to CSV
// @Description Export all article feedback data to CSV format (Admin only)
// @Tags feedback
// @Accept json
// @Produce text/csv
// @Param start_date query string false "Start date (YYYY-MM-DD)"
// @Param end_date query string false "End date (YYYY-MM-DD)"
// @Security BearerAuth
// @Success 200 {file} csv
// @Failure 500 {object} map[string]interface{}
// @Router /feedback/export/csv [get]
func (fc *FeedbackController) ExportFeedbackCSV(c *gin.Context) {
	var startDate, endDate *time.Time

	if start := c.Query("start_date"); start != "" {
		if t, err := time.Parse("2006-01-02", start); err == nil {
			startDate = &t
		}
	}

	if end := c.Query("end_date"); end != "" {
		if t, err := time.Parse("2006-01-02", end); err == nil {
			endDate = &t
		}
	}

	feedback, err := fc.feedbackService.GetAllFeedbackForCSV(startDate, endDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	// Set CSV headers
	filename := fmt.Sprintf("feedback_export_%s.csv", time.Now().Format("2006-01-02"))
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	c.Header("Content-Type", "text/csv")

	// Write CSV
	writer := csv.NewWriter(c.Writer)
	defer writer.Flush()

	// Write header
	header := []string{"feedback_id", "user_id", "username", "news_id", "article_title", "value", "value_text", "created_at", "updated_at"}
	writer.Write(header)

	// Write data
	for _, fb := range feedback {
		record := []string{
			strconv.Itoa(int(fb.FeedbackID)),
			strconv.Itoa(int(fb.UserID)),
			strconv.Itoa(int(fb.NewsID)),
			strconv.Itoa(fb.Value),
			fb.CreatedAt.Format("2006-01-02 15:04:05"),
			fb.UpdatedAt.Format("2006-01-02 15:04:05"),
		}
		writer.Write(record)
	}
}

// GetFeedbackStats godoc
// @Summary Get feedback statistics
// @Description Get aggregated feedback statistics for all articles (Admin only)
// @Tags feedback
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {array} dto.FeedbackStatsResponse
// @Router /feedback/stats [get]
func (fc *FeedbackController) GetFeedbackStats(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	stats, _, err := fc.feedbackService.GetFeedbackStats(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	/*
		response := gin.H{
			"data":       stats,
			"pagination": utils.GetPaginationInfo(page, limit, int(total)),
		}
	*/

	c.JSON(http.StatusOK, stats)
}

// GetAllFeedback godoc
// @Summary Get all feedback
// @Description Get all feedback with pagination (Admin only)
// @Tags feedback
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Security BearerAuth
// @Success 200 {array} models.Feedback
// @Router /feedback/all [get]
func (fc *FeedbackController) GetAllFeedback(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	feedback, _, err := fc.feedbackService.GetAllFeedback(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	/*
		response := gin.H{
			"data":       feedback,
			"pagination": utils.GetPaginationInfo(page, limit, int(total)),
		}
	*/

	c.JSON(http.StatusOK, feedback)
}
