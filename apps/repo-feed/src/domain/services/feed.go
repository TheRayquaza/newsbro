package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/redis/go-redis/v9"
	"repo_feed/src/api/dto"
	"repo_feed/src/converters"
	"repo_feed/src/data/models"
)

// newScore = originalScore * (0.5 ^ (age / halfLife))
func calculateDecayedScore(originalScore float64, publishedAt time.Time, now time.Time, halfLife time.Duration) float64 {
	if originalScore <= 0 {
		return 0
	}
	age := now.Sub(publishedAt)
	if age <= 0 || halfLife <= 0 {
		return originalScore
	}
	decayFactor := math.Pow(0.5, age.Hours()/halfLife.Hours())
	newScore := originalScore * decayFactor
	return newScore
}

type FeedService struct {
	RDB                *redis.Client
	defaultModel       string
	feedbackExpiration time.Duration
	decayHalfLife      time.Duration
}

func NewFeedService(rdb *redis.Client, defaultModel string, feedbackExpiration time.Duration, decayHalfLife time.Duration) *FeedService {
	return &FeedService{
		RDB:                rdb,
		defaultModel:       defaultModel,
		feedbackExpiration: feedbackExpiration,
		decayHalfLife:      decayHalfLife,
	}
}

func (s *FeedService) GetUserFeed(userID uint, model string, limit int64) ([]models.ArticleModel, error) {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}
	key := fmt.Sprintf("feed:%d:%s", userID, model)
	log.Println("Fetching feed for user ID:", userID, "with model:", model)

	zItems, err := s.RDB.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve feed ZSET: %w", err)
	}

	if len(zItems) == 0 {
		return []models.ArticleModel{}, nil
	}

	contentKeys := make([]string, 0, len(zItems))
	articleDecayScores := make(map[uint]float64)

	for _, item := range zItems {
		memberStr, ok := item.Member.(string)
		if !ok {
			log.Printf("Non-string member found in ZSET %s: %v. Skipping.", key, item.Member)
			continue
		}

		articleID, _, _, err := decodeZSetMember(memberStr)
		if err != nil {
			log.Printf("Failed to decode article ID from ZSET member '%s': %v", memberStr, err)
			continue
		}

		id := articleID
		contentKeys = append(contentKeys, fmt.Sprintf("article:%d", id))
		articleDecayScores[id] = item.Score
	}

	contents, err := s.RDB.MGet(ctx, contentKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve article contents: %w", err)
	}

	articles := make([]models.ArticleModel, 0, len(contents))
	for i, content := range contents {
		if content == nil {
			log.Printf("Article content not found for key: %s. Removing from ZSET not handled here.", contentKeys[i])
			continue
		}

		var article models.ArticleModel
		jsonStr := content.(string)
		if err := json.Unmarshal([]byte(jsonStr), &article); err != nil {
			log.Printf("Failed to unmarshal article JSON from key %s: %v", contentKeys[i], err)
			continue
		}

		article.DecayScore = articleDecayScores[article.ID]
		articles = append(articles, article)
	}

	return articles, nil
}

func (s *FeedService) UpdateFeedZSET(req *dto.UpdateFeedRequest) error {
	ctx := context.Background()
	model := req.Model
	if model == "" {
		model = s.defaultModel
	}
	feedKey := fmt.Sprintf("feed:%d:%s", req.UserID, model)
	articleContentKey := fmt.Sprintf("article:%d", req.Article.ID)

	if s.hasFeedback(ctx, req.UserID, req.Article.ID) {
		log.Printf("Skipping article ID %d for user ID %d due to negative feedback", req.Article.ID, req.UserID)
		return nil
	}

	articleModel := converters.ArticleToArticleModel(req.Article, req.Score)

	decayScore := calculateDecayedScore(articleModel.Score, req.Article.PublishedAt, time.Now(), s.decayHalfLife)
	articleModel.DecayScore = decayScore

	data, err := json.Marshal(articleModel)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	pipe := s.RDB.Pipeline()
	pipe.Set(ctx, articleContentKey, data, 0)
	memberStr := encodeZSetMember(articleModel.ID, articleModel.Score, articleModel.PublishedAt)
	pipe.ZAdd(ctx, feedKey, redis.Z{
		Score:  articleModel.DecayScore,
		Member: memberStr,
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute Redis pipeline for ZSET update and content storage: %w", err)
	}

	return nil
}

/*
func (s *FeedService) RemoveArticleZSET(userID, articleID uint, model string) error {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}
	key := fmt.Sprintf("feed:%d:%s", userID, model)

	member := encodeZSetMember(articleID, 0, time.Time{})

	if err := s.RDB.ZRem(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to remove article ID %d from ZSET %s: %w", articleID, key, err)
	}

	log.Printf("Successfully removed article ID %d from feed %s", articleID, key)
	return nil
}
*/

func (s *FeedService) AddNewFeedback(userID uint, articleID uint) error {
	ctx := context.Background()
	key := fmt.Sprintf("feedback:%d", userID)
	log.Printf("Adding negative feedback for user %d on article %d", userID, articleID)

	// Add article ID to user's negative feedback set
	if err := s.RDB.SAdd(ctx, key, articleID).Err(); err != nil {
		log.Println("Failed to add feedback to Redis:", err)
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	// Set expiration on feedback
	if s.feedbackExpiration > 0 {
		if err := s.RDB.Expire(ctx, key, s.feedbackExpiration).Err(); err != nil {
			log.Println("Failed to set expiration on feedback:", err)
		}
	}
	log.Printf("Successfully saved feedback for user %d", userID)
	return nil
}

func (s *FeedService) hasFeedback(ctx context.Context, userID uint, articleID uint) bool {
	key := fmt.Sprintf("feedback:%d", userID)
	isMember, err := s.RDB.SIsMember(ctx, key, articleID).Result()
	if err != nil {
		log.Printf("Error checking feedback for user %d, article %d: %v", userID, articleID, err)
		return false
	}
	return isMember
}
