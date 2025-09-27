package services

import (
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
	"repo_article/src/domain/entities"
)

type FeedbackService struct {
	Db *gorm.DB
}

func NewFeedbackService(db *gorm.DB) *FeedbackService {
	return &FeedbackService{
		Db: db,
	}
}

// GetArticleFeedback retrieves feedback statistics for a specific article and user's feedback
func (fs *FeedbackService) GetArticleFeedback(newsID, userID uint) (*dto.FeedbackStatsResponse, *models.Feedback, error) {
	var stats dto.FeedbackStatsResponse
	var userFeedback *models.Feedback

	// Get article feedback statistics
	err := fs.Db.Raw(`
		SELECT 
			? as news_id,
			COALESCE(SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END), 0) as like_count,
			COALESCE(SUM(CASE WHEN value = 0 THEN 1 ELSE 0 END), 0) as dislike_count,
			COALESCE(COUNT(*), 0) as total_count,
			CASE 
				WHEN COUNT(*) > 0 THEN ROUND((SUM(CASE WHEN value = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(*)), 2)
				ELSE 0 
			END as like_ratio
		FROM feedbacks 
		WHERE news_id = ? AND deleted_at IS NULL
	`, newsID, newsID).Scan(&stats).Error

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get feedback statistics: %w", err)
	}

	// Get user's specific feedback for this article
	var feedback models.Feedback
	err = fs.Db.Where("user_id = ? AND news_id = ?", userID, newsID).First(&feedback).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil, fmt.Errorf("failed to get user feedback: %w", err)
	}

	if err == nil {
		userFeedback = &models.Feedback{
			ID:        feedback.ID,
			UserID:    feedback.UserID,
			NewsID:    feedback.NewsID,
			Value:     feedback.Value,
			CreatedAt: feedback.CreatedAt,
			UpdatedAt: feedback.UpdatedAt,
		}
	}

	return &stats, userFeedback, nil
}

// CreateOrUpdateFeedback creates new feedback or updates existing one
func (fs *FeedbackService) CreateOrUpdateFeedback(userID, newsID uint, value int) (*models.Feedback, error) {
	// Validate that the article exists
	var article models.Article
	if err := fs.Db.First(&article, newsID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("article not found")
		}
		return nil, fmt.Errorf("failed to verify article existence: %w", err)
	}

	// Check if feedback already exists
	var feedback models.Feedback
	err := fs.Db.Where("user_id = ? AND news_id = ?", userID, newsID).First(&feedback).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing feedback: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new feedback
		feedback = models.Feedback{
			UserID: userID,
			NewsID: newsID,
			Value:  value,
		}

		if err := fs.Db.Create(&feedback).Error; err != nil {
			return nil, fmt.Errorf("failed to create feedback: %w", err)
		}
	} else {
		// Update existing feedback
		feedback.Value = value
		if err := fs.Db.Save(&feedback).Error; err != nil {
			return nil, fmt.Errorf("failed to update feedback: %w", err)
		}
	}

	return &models.Feedback{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		NewsID:    feedback.NewsID,
		Value:     feedback.Value,
		CreatedAt: feedback.CreatedAt,
		UpdatedAt: feedback.UpdatedAt,
	}, nil
}

// UpdateFeedback updates existing feedback
func (fs *FeedbackService) UpdateFeedback(userID, newsID uint, value int) (*models.Feedback, error) {
	var feedback models.Feedback
	err := fs.Db.Where("user_id = ? AND news_id = ?", userID, newsID).First(&feedback).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("feedback not found")
		}
		return nil, fmt.Errorf("failed to find feedback: %w", err)
	}

	feedback.Value = value
	if err := fs.Db.Save(&feedback).Error; err != nil {
		return nil, fmt.Errorf("failed to update feedback: %w", err)
	}

	return &models.Feedback{
		ID:        feedback.ID,
		UserID:    feedback.UserID,
		NewsID:    feedback.NewsID,
		Value:     feedback.Value,
		CreatedAt: feedback.CreatedAt,
		UpdatedAt: feedback.UpdatedAt,
	}, nil
}

// DeleteFeedback removes user's feedback for an article
func (fs *FeedbackService) DeleteFeedback(userID, newsID uint) error {
	result := fs.Db.Where("user_id = ? AND news_id = ?", userID, newsID).Delete(&models.Feedback{})

	if result.Error != nil {
		return fmt.Errorf("failed to delete feedback: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("feedback not found")
	}

	return nil
}

// GetUserFeedback retrieves all feedback given by a specific user with pagination
func (fs *FeedbackService) GetUserFeedback(userID uint, page, limit int) ([]models.Feedback, int64, error) {
	var feedback []models.Feedback
	var total int64

	// Count total records
	if err := fs.Db.Model(&models.Feedback{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count user feedback: %w", err)
	}

	// Get paginated results with article information
	offset := (page - 1) * limit
	err := fs.Db.Preload("Article").
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&feedback).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user feedback: %w", err)
	}

	// Convert to response format
	var response []models.Feedback
	for _, fb := range feedback {
		response = append(response, models.Feedback{
			ID:        fb.ID,
			UserID:    fb.UserID,
			NewsID:    fb.NewsID,
			Value:     fb.Value,
			CreatedAt: fb.CreatedAt,
			UpdatedAt: fb.UpdatedAt,
		})
	}

	return response, total, nil
}

// GetAllFeedbackForCSV retrieves all feedback data for CSV export
func (fs *FeedbackService) GetAllFeedbackForCSV(startDate, endDate *time.Time) ([]entities.FeedbackCSVExport, error) {
	query := fs.Db.Table("feedbacks f").
		Select(`
			f.id as feedback_id,
			f.user_id,
			f.news_id,
			f.value,
			CASE 
				WHEN f.value = 1 THEN 'Like'
				WHEN f.value = 0 THEN 'Dislike'
				ELSE 'Unknown'
			END as value_text,
			f.created_at,
			f.updated_at
		`).
		Joins("LEFT JOIN users u ON f.user_id = u.id").
		Joins("LEFT JOIN articles a ON f.news_id = a.id").
		Where("f.deleted_at IS NULL")

	if startDate != nil {
		query = query.Where("f.created_at >= ?", *startDate)
	}
	if endDate != nil {
		endTime := endDate.Add(24*time.Hour - time.Second)
		query = query.Where("f.created_at <= ?", endTime)
	}

	var feedback []entities.FeedbackCSVExport
	if err := query.Order("f.created_at DESC").Find(&feedback).Error; err != nil {
		return nil, fmt.Errorf("failed to get feedback for CSV export: %w", err)
	}

	return feedback, nil
}

// GetFeedbackStats retrieves aggregated feedback statistics for all articles
func (fs *FeedbackService) GetFeedbackStats(page, limit int) ([]dto.FeedbackStatsResponse, int64, error) {
	var stats []dto.FeedbackStatsResponse
	var total int64

	// Count total articles with feedback
	if err := fs.Db.Raw(`
		SELECT COUNT(DISTINCT news_id) 
		FROM feedbacks 
		WHERE deleted_at IS NULL
	`).Scan(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count feedback stats: %w", err)
	}

	// Get paginated statistics
	offset := (page - 1) * limit
	err := fs.Db.Raw(`
		SELECT 
			f.news_id,
			a.title as article_title,
			COALESCE(SUM(CASE WHEN f.value = 1 THEN 1 ELSE 0 END), 0) as like_count,
			COALESCE(SUM(CASE WHEN f.value = 0 THEN 1 ELSE 0 END), 0) as dislike_count,
			COALESCE(COUNT(f.id), 0) as total_count,
			CASE 
				WHEN COUNT(f.id) > 0 THEN ROUND((SUM(CASE WHEN f.value = 1 THEN 1 ELSE 0 END) * 100.0 / COUNT(f.id)), 2)
				ELSE 0 
			END as like_ratio
		FROM feedbacks f
		LEFT JOIN articles a ON f.news_id = a.id
		WHERE f.deleted_at IS NULL
		GROUP BY f.news_id, a.title
		ORDER BY total_count DESC, like_ratio DESC
		LIMIT ? OFFSET ?
	`, limit, offset).Scan(&stats).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get feedback statistics: %w", err)
	}

	return stats, total, nil
}

// GetAllFeedback retrieves all feedback with pagination (for admin)
func (fs *FeedbackService) GetAllFeedback(page, limit int) ([]models.Feedback, int64, error) {
	var feedback []models.Feedback
	var total int64

	// Count total records
	if err := fs.Db.Model(&models.Feedback{}).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count all feedback: %w", err)
	}

	// Get paginated results with user and article information
	offset := (page - 1) * limit
	err := fs.Db.Preload("User").
		Preload("Article").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&feedback).Error

	if err != nil {
		return nil, 0, fmt.Errorf("failed to get all feedback: %w", err)
	}

	return feedback, total, nil
}

// GetArticleFeedbackCount gets the total feedback count for a specific article
func (fs *FeedbackService) GetArticleFeedbackCount(newsID uint) (int64, int64, error) {
	var likeCount, dislikeCount int64

	err := fs.Db.Model(&models.Feedback{}).
		Where("news_id = ? AND value = 1", newsID).
		Count(&likeCount).Error
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count likes: %w", err)
	}

	err = fs.Db.Model(&models.Feedback{}).
		Where("news_id = ? AND value = 0", newsID).
		Count(&dislikeCount).Error
	if err != nil {
		return 0, 0, fmt.Errorf("failed to count dislikes: %w", err)
	}

	return likeCount, dislikeCount, nil
}

// HasUserFeedback checks if a user has already given feedback for an article
func (fs *FeedbackService) HasUserFeedback(userID, newsID uint) (bool, *models.Feedback, error) {
	var feedback models.Feedback
	err := fs.Db.Where("user_id = ? AND news_id = ?", userID, newsID).First(&feedback).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil, nil
		}
		return false, nil, fmt.Errorf("failed to check user feedback: %w", err)
	}

	return true, &feedback, nil
}
