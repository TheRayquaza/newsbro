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
	rdb                *redis.Client
	defaultModel       string
	feedbackExpiration time.Duration
}

func NewFeedService(rdb *redis.Client, defaultModel string, feedbackExpiration time.Duration) *FeedService {
	return &FeedService{
		rdb:                rdb,
		defaultModel:       defaultModel,
		feedbackExpiration: feedbackExpiration,
	}
}

func (s *FeedService) GetUserFeed(userId uint, model string) (*models.FeedModel, error) {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}

	key := fmt.Sprintf("feed:%d:%s", userId, model)
	log.Println("Fetching feed for user ID:", userId, "with model:", model)

	data, err := s.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		log.Println("Feed not found in Redis for key:", key)
		return nil, &dto.NotFoundError{Message: "Feed not found"}
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
		// Filter out articles that user has given feedback on the given article
		if !s.hasFeedback(ctx, req.UserID, req.Article.ID) {
			feed.Articles = append(feed.Articles, converters.ArticleToArticleModel(req.Article))
		} else {
			log.Printf("Skipping article %d due to negative feedback from user %d", req.Article.ID, req.UserID)
		}
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

func (s *FeedService) AddNewFeedback(userID uint, articleID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("feedback:%d", userID)

	log.Printf("Adding negative feedback for user %d on article %d", userID, articleID)

	// Add article ID to user's negative feedback set
	if err := s.rdb.SAdd(ctx, key, articleID).Err(); err != nil {
		log.Println("Failed to add feedback to Redis:", err)
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	// Set expiration on feedback
	if s.feedbackExpiration > 0 {
		if err := s.rdb.Expire(ctx, key, 90*24*time.Hour).Err(); err != nil {
			log.Println("Failed to set expiration on feedback:", err)
		}
	}

	log.Printf("Successfully saved feedback for user %d", userID)
	return nil
}

func (s *FeedService) hasFeedback(ctx context.Context, userID uint, articleID uint) bool {
	key := fmt.Sprintf("feedback:%d", userID)

	isMember, err := s.rdb.SIsMember(ctx, key, articleID).Result()
	if err != nil {
		log.Printf("Error checking feedback for user %d, article %d: %v", userID, articleID, err)
		return false
	}

	return isMember
}
