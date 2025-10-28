package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"

	"repo_feed/src/api/dto"
	"repo_feed/src/converters"
	"repo_feed/src/data/models"
)

type FeedService struct {
	rdb          *redis.Client
	defaultModel string
}

func NewFeedService(rdb *redis.Client, defaultModel string) *FeedService {
	return &FeedService{
		rdb:          rdb,
		defaultModel: defaultModel,
	}
}

func (s *FeedService) GetFeed(id uint, model string) (*models.FeedModel, error) {
	ctx := context.Background()

	if model == "" {
		model = s.defaultModel
	}

	key := fmt.Sprintf("feed:%d:%s", id, model)
	log.Println("Fetching feed for user ID:", id, "with model:", model)

	data, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Feed not found in Redis for key:", key)
		return nil, fmt.Errorf("feed not found")
	} else if err != nil {
		log.Println("Redis error:", err)
		return nil, fmt.Errorf("failed to retrieve feed: %w", err)
	}

	var feed models.FeedModel
	if err := json.Unmarshal([]byte(data), &feed); err != nil {
		log.Println("Failed to unmarshal feed data:", err)
		return nil, fmt.Errorf("failed to parse feed data: %w", err)
	}

	return &feed, nil
}

func (s *FeedService) UpdateFeed(req *dto.UpdateFeedRequest, add bool) error {
	ctx := context.Background()

	model := req.Model
	if model == "" {
		model = s.defaultModel
	}

	key := fmt.Sprintf("feed:%d:%s", req.UserID, model)
	log.Println("Updating feed for user ID:", req.UserID, "with model:", model)

	var feed models.FeedModel
	data, err := s.rdb.Get(ctx, key).Result()
	if err != nil && err != redis.Nil {
		log.Println("Redis error:", err)
		return fmt.Errorf("failed to retrieve feed: %w", err)
	}
	if err == nil {
		if err := json.Unmarshal([]byte(data), &feed); err != nil {
			log.Println("Failed to unmarshal feed data:", err)
			return fmt.Errorf("failed to parse feed data: %w", err)
		}
	} else {
		feed = models.FeedModel{
			UserID:    req.UserID,
			ModelName: model,
			Articles:  []models.ArticleModel{},
			UpdatedAt: time.Now(),
		}
	}
	if add {
		feed.Articles = append(feed.Articles, converters.ArticleToArticleModel(req.Article))
	} else {
		for i, article := range feed.Articles {
			if article.ID == req.Article.ID {
				feed.Articles = append(feed.Articles[:i], feed.Articles[i+1:]...)
				break
			}
		}
	}
	feed.UpdatedAt = time.Now()
	updatedData, err := json.Marshal(feed)
	if err != nil {
		log.Println("Failed to marshal updated feed data:", err)
		return fmt.Errorf("failed to serialize feed data: %w", err)
	}
	if err := s.rdb.Set(ctx, key, updatedData, 0).Err(); err != nil {
		log.Println("Failed to update feed in Redis:", err)
		return fmt.Errorf("failed to update feed: %w", err)
	}
	return nil
}
