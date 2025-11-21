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
	scoreKey           string
	feedbackKey        string
	feedKey            string
	articleKey         string
	feedbackExpiration time.Duration
	articleExpiration  time.Duration
	scoreExpiration    time.Duration
	decayHalfLife      time.Duration
	models             []string
}

func NewFeedService(rdb *redis.Client, defaultModel string, scoreKey string, feedbackKey string, feedKey string, articleKey string, feedbackExpiration time.Duration, articleExpiration time.Duration, scoreExpiration time.Duration, decayHalfLife time.Duration, models []string) *FeedService {
	return &FeedService{
		RDB:                rdb,
		defaultModel:       defaultModel,
		scoreKey:           scoreKey,
		feedbackKey:        feedbackKey,
		feedKey:            feedKey,
		articleKey:         articleKey,
		feedbackExpiration: feedbackExpiration,
		articleExpiration:  articleExpiration,
		scoreExpiration:    scoreExpiration,
		decayHalfLife:      decayHalfLife,
		models:             models,
	}
}

func (s *FeedService) GetModels() []string {
	return s.models
}

func (s *FeedService) RemoveArticleFromFeed(userID, articleID uint, model string) error {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}

	key := fmt.Sprintf("%s:%d:%s", s.feedKey, userID, model)

	data, err := s.RDB.Get(ctx, fmt.Sprintf("%s:%d", s.articleKey, articleID)).Result()
	if err == redis.Nil {
		log.Printf("article %d not found in Redis, but continuing with feed removal for user %d and model %s", articleID, userID, model)
	} else if err != nil {
		return fmt.Errorf("failed to fetch article %d when removing from feed model %s with user %d: %w", articleID, model, userID, err)
	}

	if err == nil && data != "" {
		var article models.ArticleModel
		if err := json.Unmarshal([]byte(data), &article); err != nil {
			log.Printf("failed to unmarshal article data for ID %d: %v", articleID, err)
		} else {
			scoreKey := fmt.Sprintf("%s:%d:%d:%s", s.scoreKey, userID, articleID, model)
			score := 0.0
			if scoreStr, err := s.RDB.Get(ctx, scoreKey).Result(); err == nil {
				if err := json.Unmarshal([]byte(scoreStr), &score); err != nil {
					log.Printf("failed to unmarshal score: %v", err)
				}
			}

			member := encodeZSetMember(article.ID, score, article.PublishedAt)
			if err := s.RDB.ZRem(ctx, key, member).Err(); err != nil {
				log.Printf("failed to remove article ID %d from ZSET %s: %v", articleID, key, err)
			} else {
				log.Printf("successfully removed article ID %d from feed %s", articleID, key)
				return nil
			}
		}
	}

	log.Printf("article ID %d not found in feed %s when trying to remove", articleID, key)
	return nil
}

func (s *FeedService) GetUserFeed(userID uint, model string, limit int64) ([]models.ArticleModel, error) {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}
	key := fmt.Sprintf("%s:%d:%s", s.feedKey, userID, model)

	zItems, err := s.RDB.ZRevRangeWithScores(ctx, key, 0, limit-1).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve feed ZSET for user %d and model %s: %w", userID, model, err)
	}

	if len(zItems) == 0 {
		return []models.ArticleModel{}, nil
	}

	contentKeys := make([]string, 0, len(zItems))
	scoreKeys := make([]string, 0, len(zItems))
	articleIDs := make([]uint, 0, len(zItems))
	articleDecayScores := make(map[uint]float64)

	for _, item := range zItems {
		memberStr, ok := item.Member.(string)
		if !ok {
			log.Printf("non-string member found in ZSET %s: %v. Skipping.", key, item.Member)
			continue
		}

		articleID, _, _, err := decodeZSetMember(memberStr)
		if err != nil {
			log.Printf("failed to decode article ID from ZSET member '%s': %v", memberStr, err)
			continue
		}

		articleIDs = append(articleIDs, articleID)
		contentKeys = append(contentKeys, fmt.Sprintf("%s:%d", s.articleKey, articleID))
		scoreKeys = append(scoreKeys, fmt.Sprintf("%s:%d:%d:%s", s.scoreKey, userID, articleID, model))
		articleDecayScores[articleID] = item.Score
	}

	// Fetch article contents
	contents, err := s.RDB.MGet(ctx, contentKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve article contents for user %d and model %s: %w", userID, model, err)
	}

	// Fetch user-specific scores
	scores, err := s.RDB.MGet(ctx, scoreKeys...).Result()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve scores for user %d and model %s: %w", userID, model, err)
	}

	articles := make([]models.ArticleModel, 0, len(contents))
	for i, content := range contents {
		if content == nil {
			log.Printf("article content not found for key: %s", contentKeys[i])
			continue
		}

		var article models.ArticleModel
		jsonStr := content.(string)
		if err := json.Unmarshal([]byte(jsonStr), &article); err != nil {
			log.Printf("skipping after unmarshal article JSON from key %s: %v", contentKeys[i], err)
			continue
		}

		article.Score = scores[i].(float64)

		article.DecayScore = articleDecayScores[articleIDs[i]]
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
	feedKey := fmt.Sprintf("%s:%d:%s", s.feedKey, req.UserID, model)
	articleContentKey := fmt.Sprintf("%s:%d", s.articleKey, req.Article.ID)
	scoreKey := fmt.Sprintf("%s:%d:%d:%s", s.scoreKey, req.UserID, req.Article.ID, model)

	if s.hasFeedback(ctx, req.UserID, req.Article.ID) {
		log.Printf("skipping article ID %d for user ID %d due to feedback", req.Article.ID, req.UserID)
		return nil
	}

	articleModel := converters.ArticleToArticleModel(req.Article)

	decayScore := calculateDecayedScore(req.Score, req.Article.PublishedAt, time.Now(), s.decayHalfLife)

	data, err := json.Marshal(articleModel)
	if err != nil {
		return fmt.Errorf("failed to marshal article: %w", err)
	}

	scoreData, err := json.Marshal(req.Score)
	if err != nil {
		return fmt.Errorf("failed to marshal score: %w", err)
	}

	pipe := s.RDB.Pipeline()

	pipe.Set(ctx, articleContentKey, data, s.articleExpiration)
	pipe.Set(ctx, scoreKey, scoreData, s.scoreExpiration)

	memberStr := encodeZSetMember(articleModel.ID, req.Score, articleModel.PublishedAt)
	pipe.ZAdd(ctx, feedKey, redis.Z{
		Score:  decayScore,
		Member: memberStr,
	})

	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("failed to execute Redis pipeline for ZSET update and content storage: %w", err)
	}

	return nil
}

func (s *FeedService) AddNewFeedback(userID uint, articleID uint, model string) error {
	ctx := context.Background()
	if model == "" {
		model = s.defaultModel
	}

	feedbackKey := fmt.Sprintf("%s:%d", s.feedbackKey, userID)
	log.Printf("Adding feedback for user %d on article %d", userID, articleID)

	pipe := s.RDB.Pipeline()

	pipe.SAdd(ctx, feedbackKey, articleID)

	if s.feedbackExpiration > 0 {
		pipe.Expire(ctx, feedbackKey, s.feedbackExpiration)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("Failed to add feedback to Redis:", err)
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	log.Printf("Successfully saved feedback for user %d on article %d", userID, articleID)

	if err := s.RemoveArticleFromFeed(userID, articleID, model); err != nil {
		log.Printf("Failed to remove article %d from feed for user %d: %v", articleID, userID, err)
	}

	return nil
}

func (s *FeedService) AddNewFeedbackAllModels(userID uint, articleID uint) error {
	ctx := context.Background()
	feedbackKey := fmt.Sprintf("%s:%d", s.feedbackKey, userID)
	log.Printf("Adding feedback for user %d on article %d (all models)", userID, articleID)

	for _, model := range s.models {
		if err := s.RemoveArticleFromFeed(userID, articleID, model); err != nil {
			log.Printf("Failed to remove article %d from feed %s for user %d: %v", articleID, model, userID, err)
		}
	}

	pipe := s.RDB.Pipeline()

	// Add feedback
	pipe.SAdd(ctx, feedbackKey, articleID)

	if s.feedbackExpiration > 0 {
		pipe.Expire(ctx, feedbackKey, s.feedbackExpiration)
	}

	// Remove user score
	for _, model := range s.models {
		scoreKey := fmt.Sprintf("%s:%d:%d:%s", s.scoreKey, userID, articleID, model)
		pipe.Del(ctx, scoreKey)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		log.Println("Failed to add feedback to Redis:", err)
		return fmt.Errorf("failed to save feedback: %w", err)
	}

	log.Printf("Successfully saved feedback and removed score for user %d on article %d", userID, articleID)

	return nil
}

func (s *FeedService) hasFeedback(ctx context.Context, userID uint, articleID uint) bool {
	key := fmt.Sprintf("%s:%d", s.feedbackKey, userID)
	isMember, err := s.RDB.SIsMember(ctx, key, articleID).Result()
	if err != nil {
		log.Printf("Error checking feedback for user %d, article %d: %v", userID, articleID, err)
		return false
	}
	return isMember
}
