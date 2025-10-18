package main

import (
	"context"
	"github.com/IBM/sarama"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"srvc_scrapping/src/api/consumer"
	"srvc_scrapping/src/config"
	"srvc_scrapping/src/data/models"
	"srvc_scrapping/src/data/repository"
	"srvc_scrapping/src/domain/service"
	"time"
)

func main() {
	cfg := config.Load()
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	producer, err := initKafkaProducer(cfg)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer producer.Close()
	articleRepo := repository.NewArticleRepository(db)
	rssRepo := repository.NewRSSRepository(db)
	rssService := service.NewRSSService(articleRepo, rssRepo, producer, cfg)
	log.Println("Starting RSS consumer to catch up with existing messages...")
	rssConsumer, err := consumer.NewRSSConsumer(
		cfg.KafkaBrokers,
		cfg.KafkaRSSAggregateTopic,
		cfg.KafkaGroupID,
		rssService,
	)
	if err != nil {
		log.Fatalf("Failed to create RSS consumer: %v", err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("Consuming existing RSS messages...")
	if err := rssConsumer.ConsumeUntilEmpty(ctx, 10*time.Second); err != nil {
		log.Printf("Error during initial consume: %v", err)
	}
	log.Println("Finished consuming existing messages")

	log.Println("Starting feed processing...")
	count, err := rssService.ProcessFeed(ctx)
	if err != nil {
		log.Fatalf("Failed to process feed: %v", err)
	}
	log.Printf("Processed %d new articles", count)

	if err := rssConsumer.Close(); err != nil {
		log.Printf("Error closing consumer: %v", err)
	}

	log.Println("Processing complete. Exiting...")
}

func initDB(cfg *config.Config) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&models.Article{}, &models.RSS{}); err != nil {
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
