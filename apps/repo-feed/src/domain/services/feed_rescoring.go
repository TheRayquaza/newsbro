package services

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/TheRayquaza/newsbro/apps/libs/utils"

	"github.com/redis/go-redis/v9"
)

type RescoringConfig struct {
	Enabled           bool
	Interval          time.Duration
	BatchSize         int
	ConcurrentWorkers int
	ScoreThreshold    float64
	DecayEnabled      bool
	DecayHalfLife     time.Duration
}

type FeedRescoringService struct {
	feedService *FeedService
	config      RescoringConfig
	ticker      *time.Ticker
	quit        chan struct{}
	cursorMu    sync.Mutex
}

func NewFeedRescoringService(feedService *FeedService, config RescoringConfig) *FeedRescoringService {
	return &FeedRescoringService{
		feedService: feedService,
		config:      config,
		quit:        make(chan struct{}),
	}
}

func (j *FeedRescoringService) Start() {
	if !j.config.Enabled {
		utils.SugarLog.Warn("Feed rescoring job is disabled")
		return
	}

	j.ticker = time.NewTicker(j.config.Interval)
	utils.SugarLog.Infof("Starting feed rescoring job - interval: %v, batch: %d, workers: %d, threshold: %.2f",
		j.config.Interval, j.config.BatchSize, j.config.ConcurrentWorkers, j.config.ScoreThreshold)

	go func() {
		for {
			select {
			case <-j.ticker.C:
				startTime := time.Now()
				processed, err := j.rescoreAllFeeds()
				duration := time.Since(startTime)

				if err != nil {
					utils.SugarLog.Errorf("Error during rescoring batch: %v", err)
				} else {
					utils.SugarLog.Infof("Processed %d feeds in %v (%.0f feeds/sec)",
						processed, duration, float64(processed)/duration.Seconds())
				}

			case <-j.quit:
				utils.SugarLog.Info("Stopping feed rescoring job")
				j.ticker.Stop()
				return
			}
		}
	}()
}

func (j *FeedRescoringService) Stop() {
	close(j.quit)
}

func (j *FeedRescoringService) rescoreAllFeeds() (int, error) {
	ctx := context.Background()
	var cursor uint64
	totalProcessed := 0

	for {
		keys, newCursor, err := j.feedService.RDB.Scan(ctx, cursor, "feed:*", int64(j.config.BatchSize)).Result()
		if err != nil {
			utils.SugarLog.Errorf("Failed to scan feed keys: %v", err)
			return totalProcessed, fmt.Errorf("failed to scan feed keys: %w", err)
		}

		if len(keys) > 0 {
			if err := j.processFeedsParallel(ctx, keys); err != nil {
				utils.SugarLog.Errorf("Error processing feeds batch: %v", err)
			}
			totalProcessed += len(keys)
		}

		cursor = newCursor
		if cursor == 0 {
			break
		}
	}

	return totalProcessed, nil
}

func (j *FeedRescoringService) processFeedsParallel(ctx context.Context, keys []string) error {
	jobs := make(chan string, len(keys))
	var wg sync.WaitGroup

	for i := 0; i < j.config.ConcurrentWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for key := range jobs {
				if err := j.rescoreFeedZSET(ctx, key); err != nil {
					utils.SugarLog.Errorf("worker %d: failed to rescore feed %s: %v", workerID, key, err)
				}
			}
		}(i)
	}

	for _, key := range keys {
		jobs <- key
	}
	close(jobs)

	wg.Wait()
	return nil
}

func (j *FeedRescoringService) rescoreFeedZSET(ctx context.Context, key string) error {
	zItems, err := j.feedService.RDB.ZRangeWithScores(ctx, key, 0, -1).Result()
	if err != nil {
		utils.SugarLog.Errorf("Failed to read ZSET %s: %v", key, err)
		return fmt.Errorf("failed to read ZSET %s: %w", key, err)
	}

	if len(zItems) == 0 {
		return nil
	}

	now := time.Now()

	pipe := j.feedService.RDB.Pipeline()
	articlesToUpdate := make([]redis.Z, 0)

	for _, item := range zItems {
		memberStr, ok := item.Member.(string)
		if !ok {
			utils.SugarLog.Errorf("Non-string member found in ZSET %s: %v. Skipping.", key, item.Member)
			continue
		}

		articleID, originalScore, publishedAt, err := decodeZSetMember(memberStr)
		if err != nil {
			utils.SugarLog.Errorf("Failed to decode ZSET member %s for rescoring: %v. Removing ID from ZSET.", memberStr, err)
			pipe.ZRem(ctx, key, memberStr)
			continue
		}

		originalDecayScore := item.Score
		newDecayScore := originalDecayScore

		if j.config.DecayEnabled {
			newDecayScore = calculateDecayedScore(originalScore, publishedAt, now, j.config.DecayHalfLife)
		}

		if newDecayScore < j.config.ScoreThreshold {
			pipe.ZRem(ctx, key, memberStr)
			pipe.Del(ctx, fmt.Sprintf("%s:%d", j.feedService.articleKey, articleID))
		} else if newDecayScore != originalDecayScore {
			articlesToUpdate = append(articlesToUpdate, redis.Z{
				Score:  newDecayScore,
				Member: memberStr,
			})
		}
	}

	if len(articlesToUpdate) > 0 {
		pipe.ZAdd(ctx, key, articlesToUpdate...)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		utils.SugarLog.Errorf("Failed to execute rescoring pipeline for ZSET %s: %v", key, err)
		return fmt.Errorf("failed to execute rescoring pipeline for ZSET %s: %w", key, err)
	}

	return nil
}
