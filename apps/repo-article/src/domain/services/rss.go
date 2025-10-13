package services

import (
	"repo_article/src/api/dto"
	"repo_article/src/config"
	"repo_article/src/converters"
	"repo_article/src/data/models"

	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"
	"gorm.io/gorm"
)

type RSSService struct {
	Db       *gorm.DB
	producer sarama.SyncProducer
	config   *config.Config
}

func NewRSSService(db *gorm.DB, producer sarama.SyncProducer, config *config.Config) *RSSService {
	return &RSSService{
		Db:       db,
		producer: producer,
		config:   config,
	}
}
func (as *RSSService) CreateRSS(req *dto.RSSCreateRequest, userID uint) (*dto.RSSResponse, error) {
	rss := &models.RSSSource{
		Name:        req.Name,
		Description: req.Description,
		Link:        req.Link,
		Active:      true,
	}

	// Handle parents assignment before creation
	if req.Parents != nil && len(*req.Parents) > 0 {
		var parents []models.RSSSource
		if err := as.Db.Where("name IN ?", *req.Parents).Find(&parents).Error; err != nil {
			log.Printf("failed to get parent rss: %s", err)
			return nil, fmt.Errorf("failed to get parent rss: %w", err)
		}

		if len(parents) != len(*req.Parents) {
			return nil, dto.NewBadRequest("some parent rss not found")
		}
		log.Printf("Assigning parents: %+v", parents)
		rss.Parents = parents
	}

	// Create the RSS source with associations
	if err := as.Db.Create(rss).Error; err != nil {
		log.Printf("failed to create rss: %s", err)
		return nil, fmt.Errorf("failed to create rss: %w", err)
	}

	// Publish to Kafka
	rssAggregate := converters.RSSModelToRSSAggregate(rss, true)
	rssBytes, err := json.Marshal(rssAggregate)
	if err != nil {
		log.Printf("failed to marshal rss aggregate: %s", err)
		return nil, err
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(rssAggregate.Name),
		Topic: as.config.KafkaRSSAggregateTopic,
		Value: sarama.ByteEncoder(rssBytes),
	}

	_, _, err = as.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to publish rss to kafka: %s", err)
		return nil, fmt.Errorf("failed to publish rss to kafka: %w", err)
	}

	rssResponse := converters.RSSModelToRSSResponse(rss)
	return &rssResponse, nil
}

func (as *RSSService) GetRSSByName(name string) (*dto.RSSResponse, error) {
	var rss models.RSSSource

	if err := as.Db.First(&rss, name).Error; err != nil {
		log.Printf("failed to get rss: %s", err)
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFound(fmt.Sprintf("rss with name %s not found", name))
		}
		return nil, fmt.Errorf("failed to get rss: %w", err)
	}

	rssResponse := converters.RSSModelToRSSResponse(&rss)

	return &rssResponse, nil
}

func (as *RSSService) GetRSS(userID uint, filters *dto.RSSFilters) ([]dto.RSSResponse, error) {
	var rss []models.RSSSource
	query := as.Db.Model(&models.RSSSource{})

	// Filtering
	if filters.BeginDate != nil {
		query = query.Where("created_at >= ?", filters.BeginDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}

	if err := query.Order("name ASC").
		Find(&rss).Error; err != nil {
		return nil, fmt.Errorf("failed to get rss: %w", err)
	}

	rssResponses := make([]dto.RSSResponse, 0, len(rss))
	for _, r := range rss {
		resp := converters.RSSModelToRSSResponse(&r)
		rssResponses = append(rssResponses, resp)
	}

	return rssResponses, nil
}

func (as *RSSService) GetTreeRSS(userID uint, filters *dto.RSSFilters) ([]dto.TreeRSSResponse, error) {
	var allFeeds []models.RSSSource
	query := as.Db.Model(&models.RSSSource{})

	if filters.BeginDate != nil {
		query = query.Where("created_at >= ?", filters.BeginDate)
	}
	if filters.EndDate != nil {
		query = query.Where("created_at <= ?", filters.EndDate)
	}

	if err := query.Preload("Children.Children").Preload("Parents").Order("link ASC").Find(&allFeeds).Error; err != nil {
		log.Printf("failed to get rss tree: %s", err)
		return nil, fmt.Errorf("failed to get rss tree: %w", err)
	}

	if len(allFeeds) == 0 {
		return []dto.TreeRSSResponse{}, nil
	}

	feedMap := make(map[string]*models.RSSSource)
	for i := range allFeeds {
		feedMap[allFeeds[i].Name] = &allFeeds[i]
	}

	rootFeeds := make(map[string]bool)
	for _, feed := range allFeeds {
		if len(feed.Parents) == 0 {
			rootFeeds[feed.Name] = true
		}
	}

	var result []dto.TreeRSSResponse
	visited := make(map[string]bool)

	for _, feed := range allFeeds {
		if rootFeeds[feed.Name] && !visited[feed.Name] {
			treeDTO := as.buildTreeDTO(&feed, visited)
			result = append(result, treeDTO)
		}
	}

	return result, nil
}

func (as *RSSService) buildTreeDTO(feed *models.RSSSource, visited map[string]bool) dto.TreeRSSResponse {
	visited[feed.Name] = true

	treeDTO := dto.TreeRSSResponse{
		Name:        feed.Name,
		Description: feed.Description,
		Link:        feed.Link,
		Active:      feed.Active,
		CreatedAt:   feed.CreatedAt,
		UpdatedAt:   feed.UpdatedAt,
		Children:    []dto.TreeRSSResponse{},
	}

	for i := range feed.Children {
		child := feed.Children[i]
		log.Printf("Processing child: %s of parent: %s", child.Name, feed.Name)
		if !visited[child.Name] {
			childDTO := as.buildTreeDTO(&child, visited)
			treeDTO.Children = append(treeDTO.Children, childDTO)
		}
	}	
	log.Printf("Visiting: %s, %d children found", treeDTO.Name, len(treeDTO.Children))

	return treeDTO
}

func (as *RSSService) UpdateRSS(name string, req *dto.RSSUpdateRequest) (*dto.RSSResponse, error) {
	var rss models.RSSSource
	if err := as.Db.First(&rss, name).Error; err != nil {
		log.Printf("failed to get rss: %s", err)
		if err == gorm.ErrRecordNotFound {
			return nil, dto.NewNotFound(fmt.Sprintf("rss with name %s not found", name))
		}
		return nil, fmt.Errorf("failed to get rss: %w", err)
	}
	if req.Description != nil {
		rss.Description = *req.Description
	}
	if req.Link != nil {
		rss.Link = *req.Link
	}
	if req.Parents != nil {
		var parents []models.RSSSource
		if len(*req.Parents) > 0 {
			if err := as.Db.Where("name IN ?", *req.Parents).Find(&parents).Error; err != nil {
				log.Printf("failed to get parent rss: %s", err)
				return nil, fmt.Errorf("failed to get parent rss: %w", err)
			}
			if len(parents) != len(*req.Parents) {
				return nil, dto.NewBadRequest("some parent rss not found")
			}
			rss.Parents = parents
		} else {
			// Clear parents if empty list provided
			if err := as.Db.Model(&rss).Association("RSSParents").Clear(); err != nil {
				log.Printf("failed to clear parent rss: %s", err)
				return nil, fmt.Errorf("failed to clear parent rss: %w", err)
			}
		}
	}
	if err := as.Db.Save(&rss).Error; err != nil {
		log.Printf("failed to update rss: %s", err)
		return nil, fmt.Errorf("failed to update rss: %w", err)
	}
	rssAggregate := converters.RSSModelToRSSAggregate(&rss, rss.Active)
	rssBytes, err := json.Marshal(rssAggregate)
	if err != nil {
		log.Printf("failed to marshal rss aggregate: %s", err)
		return nil, err
	}
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(rssAggregate.Name),
		Topic: as.config.KafkaRSSAggregateTopic,
		Value: sarama.ByteEncoder(rssBytes),
	}
	_, _, err = as.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to publish rss update to kafka: %s", err)
		return nil, fmt.Errorf("failed to publish rss update to kafka: %w", err)
	}
	rssResponse := converters.RSSModelToRSSResponse(&rss)
	return &rssResponse, nil
}

func (as *RSSService) DeleteRSS(name string) error {
	var rss models.RSSSource
	if err := as.Db.First(&rss, name).Error; err != nil {
		log.Printf("failed to get rss: %s", err)
		if err == gorm.ErrRecordNotFound {
			return dto.NewNotFound(fmt.Sprintf("rss with name %s not found", name))
		}
		return fmt.Errorf("failed to get rss: %w", err)
	}
	if err := as.Db.Delete(&rss).Error; err != nil {
		log.Printf("failed to delete rss: %s", err)
		return fmt.Errorf("failed to delete rss: %w", err)
	}
	rssAggregate := converters.RSSModelToRSSAggregate(&rss, false)
	rssBytes, err := json.Marshal(rssAggregate)
	if err != nil {
		log.Printf("failed to marshal rss aggregate: %s", err)
		return err
	}
	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(rssAggregate.Name),
		Topic: as.config.KafkaRSSAggregateTopic,
		Value: sarama.ByteEncoder(rssBytes),
	}
	_, _, err = as.producer.SendMessage(msg)
	if err != nil {
		log.Printf("failed to publish rss deletion to kafka: %s", err)
		return fmt.Errorf("failed to publish rss deletion to kafka: %w", err)
	}
	return nil
}
