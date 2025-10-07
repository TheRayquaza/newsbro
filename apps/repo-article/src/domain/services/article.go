package services

import (
	"repo_article/src/api/dto"
	"repo_article/src/config"
	"repo_article/src/converters"
	"repo_article/src/data/models"

	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type ArticleService struct {
	Db       *gorm.DB
	producer sarama.SyncProducer
	config   *config.Config
}

func NewArticleService(db *gorm.DB, producer sarama.SyncProducer, config *config.Config) *ArticleService {
	return &ArticleService{
		Db:       db,
		producer: producer,
		config:   config,
	}
}

func (as *ArticleService) CreateArticle(req *dto.ArticleCreateRequest, userID uint) (*dto.ArticleResponse, error) {
	article := &models.Article{
		Category:    req.Category,
		Subcategory: req.Subcategory,
		Title:       req.Title,
		Abstract:    req.Abstract,
		Link:        req.Link,
		PublishedAt: req.PublishedAt,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if req.ID != 0 {
		var existing models.Article
		if err := as.Db.First(&existing, "id = ?", req.ID).Error; err == nil {
			log.Printf("article with id %d already exists", req.ID)
			return nil, dto.NewConflict(fmt.Sprintf("article with id %d already exists", req.ID))
		} else if err != gorm.ErrRecordNotFound {
			log.Printf("failed to check existing article: %s", err)
			return nil, dto.NewNotFound(fmt.Sprintf("failed to check existing article: %s", err))
		}
		article.ID = req.ID
	}

	if err := as.Db.Create(article).Error; err != nil {
		log.Printf("failed to create article: %s", err)
		return nil, fmt.Errorf("failed to create article: %w", err)
	}

	articleAggregate := converters.ArticleModelToArticleAggregate(article, true)

	articleBytes, err := json.Marshal(articleAggregate)
	if err != nil {
		log.Printf("failed to marshal article aggregate: %s", err)
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(articleAggregate.Link),
		Topic: as.config.KafkaArticleAggregateTopic,
		Value: sarama.ByteEncoder(articleBytes),
	}

	_, _, err = as.producer.SendMessage(msg)

	if err != nil {
		log.Printf("failed to publish article to kafka: %s", err)
		return nil, fmt.Errorf("failed to publish article to kafka: %w", err)
	}

	articleResponse := converters.ArticleModelToArticleResponse(article)
	return &articleResponse, nil
}

func (as *ArticleService) GetArticleByID(id uint) (*dto.ArticleResponse, error) {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		log.Printf("failed to get article: %s", err)
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("article with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

	articleResponse := converters.ArticleModelToArticleResponse(&article)
	return &articleResponse, nil
}

func (as *ArticleService) GetArticles(filters *dto.ArticleFilters) ([]dto.ArticleResponse, int64, error) {
	if filters.Limit < 0 {
		return nil, 0, dto.NewBadRequest("limit must be >= 0")
	}
	if filters.Offset < 0 {
		return nil, 0, dto.NewBadRequest("offset must be >= 0")
	}

	filters.Limit = min(filters.Limit, as.config.MaxPageSize)

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

	var articleResponses []dto.ArticleResponse
	for _, article := range articles {
		articleResponses = append(articleResponses, converters.ArticleModelToArticleResponse(&article))
	}

	return articleResponses, total, nil
}

func (as *ArticleService) UpdateArticle(id uint, req *dto.ArticleUpdateRequest, userID uint) (*dto.ArticleResponse, error) {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		log.Printf("failed to get article: %s", err)
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFound(fmt.Sprintf("article with id %d not found", id))
		}
		return nil, fmt.Errorf("failed to get article: %w", err)
	}

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
		log.Printf("failed to update article: %s", err)
		return nil, fmt.Errorf("failed to update article: %w", err)
	}

	if err := as.Db.First(&article, id).Error; err != nil {
		log.Printf("failed to reload updated article: %s", err)
		return nil, fmt.Errorf("failed to reload updated article: %w", err)
	}

	articleAggregate := converters.ArticleModelToArticleAggregate(&article, true)

	articleBytes, err := json.Marshal(articleAggregate)
	if err != nil {
		log.Printf("failed to marshal article aggregate: %s", err)
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(articleAggregate.Link),
		Topic: as.config.KafkaArticleAggregateTopic,
		Value: sarama.ByteEncoder(articleBytes),
	}

	_, _, err = as.producer.SendMessage(msg)

	if err != nil {
		log.Printf("failed to publish article to kafka: %s", err)
		return nil, fmt.Errorf("failed to publish article to kafka: %w", err)
	}

	articleResponse := converters.ArticleModelToArticleResponse(&article)
	return &articleResponse, nil
}

func (as *ArticleService) DeleteArticle(id uint, userID uint) error {
	var article models.Article
	if err := as.Db.First(&article, id).Error; err != nil {
		log.Printf("failed to get article: %s", err)
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFound(fmt.Sprintf("article with id %d not found", id))
		}
		return fmt.Errorf("failed to get article: %w", err)
	}

	if err := as.Db.Delete(&article).Error; err != nil {
		log.Printf("failed to delete article: %s", err)
		return fmt.Errorf("failed to delete article: %w", err)
	}

	articleAggregate := converters.ArticleModelToArticleAggregate(&article, false)

	articleBytes, err := json.Marshal(articleAggregate)
	if err != nil {
		log.Printf("failed to marshal article aggregate: %s", err)
		return err
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(article.Link),
		Topic: as.config.KafkaArticleAggregateTopic,
		Value: sarama.ByteEncoder(articleBytes),
	}

	_, _, err = as.producer.SendMessage(msg)

	if err != nil {
		log.Printf("failed to publish article to kafka: %s", err)
		return fmt.Errorf("failed to publish article to kafka: %w", err)
	}

	return nil
}

func (as *ArticleService) GetCategories() ([]string, error) {
	var categories []string
	if err := as.Db.Model(&models.Article{}).
		Distinct("category").
		Pluck("category", &categories).Error; err != nil {
		log.Printf("failed to get categories: %s", err)
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
		log.Printf("failed to get subcategories: %s", err)
		return nil, fmt.Errorf("failed to get subcategories: %w", err)
	}

	return subcategories, nil
}

func (as *ArticleService) TriggerArticleIngestion(beginDate, endDate time.Time) (uint, error) {
	if beginDate.After(endDate) {
		return 0, dto.NewBadRequest("begin_date must be before end_date")
	}

	var articles []models.Article
	if err := as.Db.Where("published_at BETWEEN ? AND ?", beginDate, endDate).Find(&articles).Error; err != nil {
		log.Printf("failed to get articles for ingestion: %s", err)
		return 0, fmt.Errorf("failed to get articles for ingestion: %w", err)
	}
	if len(articles) == 0 {
		log.Printf("no articles found for ingestion between %s and %s", beginDate, endDate)
		return 0, nil
	}

	count := uint(0)
	failed := make([]uint, 0)

	for _, article := range articles {
		articleAggregate := converters.ArticleModelToArticleAggregate(&article, true)
		articleBytes, err := json.Marshal(articleAggregate)
		if err != nil {
			log.Printf("failed to marshal article aggregate: %s", err)
			failed = append(failed, article.ID)
			continue
		}
		msg := &sarama.ProducerMessage{
			Key:   sarama.StringEncoder(articleAggregate.Link),
			Topic: as.config.KafkaArticleAggregateTopic,
			Value: sarama.ByteEncoder(articleBytes),
		}
		_, _, err = as.producer.SendMessage(msg)
		if err != nil {
			log.Printf("failed to publish article to kafka: %s", err)
			failed = append(failed, article.ID)
			continue
		}
		count++
	}

	if len(failed) > 0 {
		log.Printf("failed to ingest articles with IDs: %v", failed)
		for _, id := range failed {
			log.Printf(" - %d", id)
		}
	}

	return count, nil
}
