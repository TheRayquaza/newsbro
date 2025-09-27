package services

import (
	"fmt"
	"log"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
	"strings"

	"gorm.io/gorm"
)

type ArticleService struct {
	Db *gorm.DB
}

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{
		Db: db,
	}
}

// CreateArticle creates a new article and publishes to Kafka
func (as *ArticleService) CreateArticle(req *dto.ArticleCreateRequest, userID uint) (*models.Article, error) {
	article := &models.Article{
		Category:    req.Category,
		Subcategory: req.Subcategory,
		Title:       req.Title,
		Abstract:    req.Abstract,
		Link:        req.Link,
	}
	if req.ID != 0 {
		var existing models.Article
		if err := as.Db.First(&existing, "id = ?", req.ID).Error; err == nil {
			return nil, dto.NewConflict(fmt.Sprintf("article with id %d already exists", req.ID))
		} else if err != gorm.ErrRecordNotFound {
			return nil, dto.NewNotFound(fmt.Sprintf("failed to check existing article: %s", err))
		}
		article.ID = req.ID
	}

	if err := as.Db.Create(article).Error; err != nil {
		log.Printf("failed to create article: %s", err)
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	// Publish to Kafka
	//as.publishToKafka("article_created", article, userID)

	return article, nil
}

// GetArticleByID retrieves an article by its ID
func (as *ArticleService) GetArticleByID(id uint) (*models.Article, error) {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}
	return &article, nil
}

func (as *ArticleService) GetArticles(filters *dto.ArticleFilters) ([]models.Article, int64, error) {
	if filters.Limit < 0 {
		return nil, 0, dto.NewBadRequest("limit must be >= 0")
	}
	if filters.Offset < 0 {
		return nil, 0, dto.NewBadRequest("offset must be >= 0")
	}

	var articles []models.Article
	var total int64
	query := as.Db.Model(&models.Article{})

	if filters.Category != "" {
		query = query.Where("category = ?", filters.Category)
	}

	if filters.Subcategory != "" {
		query = query.Where("subcategory = ?", filters.Subcategory)
	}

	if filters.Search != "" {
		searchTerm := "%" + strings.ToLower(filters.Search) + "%"
		query = query.Where("LOWER(title) LIKE ? OR LOWER(abstract) LIKE ?", searchTerm, searchTerm)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	if err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

func (as *ArticleService) UpdateArticle(id uint, req *dto.ArticleUpdateRequest, userID uint) (*models.Article, error) {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFound(fmt.Sprintf("article with id %d not found", id))
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	// Update fields if provided
	updates := make(map[string]interface{})
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Subcategory != nil {
		updates["subcategory"] = *req.Subcategory
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Abstract != nil {
		updates["abstract"] = *req.Abstract
	}
	if req.Link != nil {
		updates["link"] = *req.Link
	}

	if err := as.Db.Model(&article).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	// Reload the article
	if err := as.Db.First(&article, id).Error; err != nil {
		return nil, fmt.Errorf("failed to reload updated article: %w", err)
	}

	// Publish to Kafka
	//as.publishToKafka("article_updated", &article, userID)

	return &article, nil
}

func (as *ArticleService) DeleteArticle(id uint, userID uint) error {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFound(fmt.Sprintf("article with id %d not found", id))
		}
		return fmt.Errorf("failed to get article: %w", err)
	}

	if err := as.Db.Delete(&article).Error; err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	// Publish to Kafka
	//as.publishToKafka("article_deleted", &article, userID)

	return nil
}

// GetCategories retrieves all unique categories
func (as *ArticleService) GetCategories() ([]string, error) {
	var categories []string
	if err := as.Db.Model(&models.Article{}).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		return nil, fmt.Errorf("failed to get categories: %w", err)
	}
	return categories, nil
}

func (as *ArticleService) GetSubcategories(category string) ([]string, error) {
	var subcategories []string
	query := as.Db.Model(&models.Article{}).Distinct("subcategory")

	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Pluck("subcategory", &subcategories).Error; err != nil {
		return nil, fmt.Errorf("failed to get subcategories: %w", err)
	}

	return subcategories, nil
}
