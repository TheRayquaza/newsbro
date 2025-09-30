package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"repo_article/src/api/routes"
	"repo_article/src/api/consumer"
	"repo_article/src/config"
	"repo_article/src/data/database"
	"repo_article/src/domain/services"

	"github.com/IBM/sarama"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize Kafka producer
	producer, err := initKafkaProducer(cfg)
	if err != nil {
		log.Fatal("Failed to create Kafka producer:", err)
	}

	// Initialize services
	articleService := services.NewArticleService(db, producer, cfg)
	feedbackService := services.NewFeedbackService(db, producer, cfg)

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
		log.Printf("Warning: Failed to initialize Kafka consumer: %v", err)
	} else {
		go consumer.Start(ctx)
		log.Println("Kafka consumer initialized successfully")
	}

	// Setup routes
	router := routes.SetupRouter(cfg, articleService, feedbackService)

	// Handle graceful shutdown
	go func() {
		sigterm := make(chan os.Signal, 1)
		signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
		<-sigterm
		log.Println("Shutting down gracefully...")
		cancel()
		if consumer != nil {
			consumer.Close()
		}
		os.Exit(0)
	}()

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

func initKafkaProducer(cfg *config.Config) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	return sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaConfig)
}
