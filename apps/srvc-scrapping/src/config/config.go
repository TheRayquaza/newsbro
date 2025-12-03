package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                     string
	Environment              string
	DatabaseURL              string
	KafkaBrokers             []string
	KafkaRSSAggregateTopic   string
	KafkaArticleCommandTopic string
	KafkaGroupID             string
	WebhookURL               string
}

func Map[T, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i := range ts {
		us[i] = f(ts[i])
	}
	return us
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "prod"),
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			getEnv("DATABASE_USERNAME", "username"),
			getEnv("DATABASE_PASSWORD", "password"),
			getEnv("DATABASE_HOST", "localhost:5432"),
			getEnv("DATABASE_NAME", "repo_account"),
		),
		KafkaBrokers:             strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaRSSAggregateTopic:   getEnv("KAFKA_RSS_AGGREGATE_TOPIC", "rss-aggregate"),
		KafkaGroupID:             getEnv("KAFKA_GROUP_ID", "srvc-scrapping-group"),
		KafkaArticleCommandTopic: getEnv("KAFKA_ARTICLE_COMMAND_TOPIC", "new-articles-command"),
		WebhookURL:               getEnv("WEBHOOK_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
