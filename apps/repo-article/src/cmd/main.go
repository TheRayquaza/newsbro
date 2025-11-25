package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"repo_article/src/api/consumer"
	"repo_article/src/api/routes"
	"repo_article/src/config"
	"repo_article/src/data/database"
	"repo_article/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/utils"

	"github.com/IBM/sarama"
)

func main() {
	// Load configuration
	cfg := config.Load()

	if err := utils.Initialize(cfg.Environment); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer utils.SugarLog.Sync()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		utils.SugarLog.Fatal("Failed to initialize database:", err)
	}

	// Initialize Kafka producer
	producer, err := initKafkaProducer(cfg)
	if err != nil {
		utils.SugarLog.Fatal("Failed to create Kafka producer:", err)
	}

	// Initialize services
	articleService := services.NewArticleService(db, producer, cfg)
	feedbackService := services.NewFeedbackService(db, producer, cfg)
	rssService := services.NewRSSService(db, producer, cfg)

	// Setup Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	consumer, err := consumer.NewArticleConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaArticleCommandTopic,
		cfg.KafkaGroupID,
		articleService,
	)
	if err != nil {
		utils.SugarLog.Errorf("Warning: Failed to initialize Kafka consumer: %v", err)
	} else {
		go consumer.Start(ctx)
		utils.SugarLog.Infof("Kafka consumer initialized successfully")
	}

	// Setup routes
	router := routes.SetupRouter(cfg, articleService, feedbackService, rssService)

	// Handle graceful shutdown
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		utils.SugarLog.Infof("Shutting down gracefully...")
		cancel()
		if consumer != nil {
			consumer.Close()
		}
		os.Exit(0)
	}()

	// Start server
	utils.SugarLog.Infof("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		utils.SugarLog.Fatal("Failed to start server:", err)
	}
}

func initKafkaProducer(cfg *config.Config) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	return sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaConfig)
}
