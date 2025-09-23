package services

import (
	"fmt"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
	"strings"

	//"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type ArticleService struct {
	Db *gorm.DB
	//kafkaProducer sarama.SyncProducer
}

/*
// KafkaMessage represents the message structure sent to Kafka
type KafkaMessage struct {
	Action    string          `json:"action"`
	Article   *models.Article `json:"article"`
	Timestamp time.Time       `json:"timestamp"`
	UserID    uint            `json:"user_id,omitempty"`
}
*/

func NewArticleService(db *gorm.DB) *ArticleService {
	return &ArticleService{
		Db: db,
		//kafkaProducer: nil,
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
			return nil, fmt.Errorf("article with id %s already exists", req.ID)
		} else if err != gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("failed to check existing article: %w", err)
		}
		article.ID = req.ID
	}

	if err := as.Db.Create(article).Error; err != nil {
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

// GetArticles retrieves articles with filtering and pagination
func (as *ArticleService) GetArticles(filters *dto.ArticleFilters) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := as.Db.Model(&models.Article{})

	// Apply filters
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

	// Count total records
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count articles: %w", err)
	}

	// Apply pagination and order
	if err := query.Order("created_at DESC").
		Limit(filters.Limit).
		Offset(filters.Offset).
		Find(&articles).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to get articles: %w", err)
	}

	return articles, total, nil
}

// UpdateArticle updates an existing article
func (as *ArticleService) UpdateArticle(id uint, req *dto.ArticleUpdateRequest, userID uint) (*models.Article, error) {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article with id %d not found", id)
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

	// Reload the article to get updated data
	if err := as.Db.First(&article, id).Error; err != nil {
		return nil, fmt.Errorf("failed to reload updated article: %w", err)
	}

	// Publish to Kafka
	//as.publishToKafka("article_updated", &article, userID)

	return &article, nil
}

// DeleteArticle soft deletes an article
func (as *ArticleService) DeleteArticle(id uint, userID uint) error {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("article with id %d not found", id)
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

// GetSubcategories retrieves subcategories for a given category
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

/*
// publishToKafka publishes article events to Kafka
func (as *ArticleService) publishToKafka(action string, article *models.Article, userID uint) {
	if as.kafkaProducer == nil {
		log.Println("Kafka producer not available, skipping message")
		return
	}

	message := KafkaMessage{
		Action:    action,
		Article:   article,
		Timestamp: time.Now(),
		UserID:    userID,
	}

	messageBytes, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal Kafka message: %v", err)
		return
	}

	kafkaMessage := &sarama.ProducerMessage{
		Topic: "articles",
		Key:   sarama.StringEncoder(fmt.Sprintf("article_%d", article.ID)),
		Value: sarama.ByteEncoder(messageBytes),
	}

	partition, offset, err := as.kafkaProducer.SendMessage(kafkaMessage)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v", err)
		return
	}

	log.Printf("Message sent to Kafka - Topic: articles, Partition: %d, Offset: %d", partition, offset)
}
*/
