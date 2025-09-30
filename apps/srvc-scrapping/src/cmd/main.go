package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"srvc_scrapping/src/config"
	"srvc_scrapping/src/data/models"
	"srvc_scrapping/src/data/repository"
	"srvc_scrapping/src/domain/service"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	// Initialize database
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize Kafka producer
	producer, err := initKafkaProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()

	// Initialize repositories
	articleRepo := repository.NewArticleRepository(db)

	// Initialize use cases
	rssService := service.NewRSSService(articleRepo, producer, cfg)

	// Setup HTTP router
	router := gin.Default()
	router.Use(gin.Recovery())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		if err := db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "healthy"})
	})

	// Readiness check
	router.GET("/ready", func(c *gin.Context) {
		if err := db.Raw("SELECT 1").Error; err != nil {
			c.JSON(500, gin.H{"status": "unhealthy", "reason": err})
			return
		}
		c.JSON(200, gin.H{"status": "ready"})
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Starting server on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Process RSS feed once at startup
	ctx := context.Background()
	count, err := rssService.ProcessFeed(ctx)
	if err != nil {
		log.Fatalf("Failed to process feed: %v", err)
	}
	log.Printf("Processed %d new articles", count)

	// Graceful shutdown
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	if err := db.AutoMigrate(&models.Article{}); err != nil {
		return nil, err
	}

	return db, nil
}

func initKafkaProducer(cfg *config.Config) (sarama.SyncProducer, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.RequiredAcks = sarama.WaitForAll
	kafkaConfig.Producer.Retry.Max = 5
	kafkaConfig.Producer.Return.Successes = true

	return sarama.NewSyncProducer(cfg.KafkaBrokers, kafkaConfig)
}
