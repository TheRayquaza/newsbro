package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"repo_feed/src/api/consumer"
	"repo_feed/src/api/routes"
	"repo_feed/src/config"
	"repo_feed/src/data/redis"
	"repo_feed/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/utils"
)

func main() {
	// Load configuration
	cfg := config.Load()

	if err := utils.Initialize(cfg.Environment); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer utils.SugarLog.Sync()

	if cfg.Environment == "dev" {
		utils.SugarLog.Infof("Running in development mode")
	}

	// Initialize Redis pool
	rdb, err := redis.Initialize(redis.RedisConfig{
		Sentinels:  cfg.RedisSentinels,
		MasterName: cfg.RedisMasterName,
		Password:   cfg.RedisPassword,
		DB:         cfg.RedisDB,
	})
	if err != nil {
		utils.SugarLog.Fatalf("Failed to initialize Redis: %v", err)
	}

	// Initialize feed service
	feedService := services.NewFeedService(rdb, cfg.DefaultModel, cfg.RedisScoreKey, cfg.RedisFeedbackKey, cfg.RedisFeedKey, cfg.RedisArticleKey, cfg.FeedbackExpiration, cfg.ArticleExpiration, cfg.ScoreExpiration, cfg.FeedRescoring.DecayHalfLife, cfg.Models)

	// Initialize and start rescoring job
	rescoringConfig := services.RescoringConfig{
		Interval:          cfg.FeedRescoring.Interval,
		BatchSize:         cfg.FeedRescoring.BatchSize,
		ConcurrentWorkers: cfg.FeedRescoring.ConcurrentWorkers,
		ScoreThreshold:    cfg.FeedRescoring.ScoreThreshold,
		DecayEnabled:      cfg.FeedRescoring.DecayEnabled,
		DecayHalfLife:     cfg.FeedRescoring.DecayHalfLife,
		Enabled:           cfg.FeedRescoring.Enabled,
	}

	rescoringService := services.NewFeedRescoringService(feedService, rescoringConfig)
	rescoringService.Start()
	utils.SugarLog.Infof("Feed rescoring service started")

	// Setup context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Setup Kafka inference consumer
	inferenceConsumer, err := consumer.NewInferenceConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaInferenceCommandTopic,
		cfg.KafkaGroupID,
		feedService,
	)
	if err != nil {
		utils.SugarLog.Errorf("Warning: Failed to initialize Kafka inference consumer: %v", err)
	} else {
		go inferenceConsumer.Start(ctx)
		utils.SugarLog.Infof("Kafka inference consumer initialized successfully")
	}

	// Setup Kafka feedback consumer
	feedbackConsumer, err := consumer.NewFeedbackConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaFeedbackTopic,
		cfg.KafkaFeedbackGroupID,
		feedService,
	)
	if err != nil {
		utils.SugarLog.Errorf("Warning: Failed to initialize Kafka feedback consumer: %v", err)
	} else {
		go feedbackConsumer.Start(ctx)
		utils.SugarLog.Infof("Kafka feedback consumer initialized successfully")
	}

	// Setup routes
	router := routes.SetupRouter(cfg, feedService, rescoringService)
	//router.SetTrustedProxies([]string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}) kube internal

	// Graceful shutdown handler
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm

		utils.SugarLog.Infof("Shutting down gracefully...")
		cancel()

		// Stop rescoring service
		rescoringService.Stop()

		// Close inference consumer
		if inferenceConsumer != nil {
			if err := inferenceConsumer.Close(); err != nil {
				utils.SugarLog.Errorf("Error closing inference consumer: %v", err)
			}
		}

		// Close feedback consumer
		if feedbackConsumer != nil {
			if err := feedbackConsumer.Close(); err != nil {
				utils.SugarLog.Errorf("Error closing feedback consumer: %v", err)
			}
		}

		os.Exit(0)
	}()

	// Start server
	utils.SugarLog.Infof("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		utils.SugarLog.Fatalf("Failed to start server: %v", err)
	}
}
