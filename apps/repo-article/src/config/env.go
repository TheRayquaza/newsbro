package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port             string
	DatabaseURL      string
	JWTSecret        string
	LoginRedirectURL string
	KafkaBrokers     []string
	KafkaArticleCommandTopic       string
	KafkaGroupID     string
	KafkaArticleAggregateTopic string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			getEnv("DATABASE_USERNAME", "username"),
			getEnv("DATABASE_PASSWORD", "password"),
			getEnv("DATABASE_HOST", "localhost:5432"),
			getEnv("DATABASE_NAME", "repo_account"),
		),
		JWTSecret:        getEnv("JWT_SECRET", ""),
		LoginRedirectURL: getEnv("LOGIN_REDIRECT_URL", "http://localhost:8080"),
		KafkaBrokers:     strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaArticleCommandTopic:       getEnv("KAFKA_ARTICLE_COMMAND_TOPIC", "new-articles-command"),
		KafkaArticleAggregateTopic:       getEnv("KAFKA_ARTICLE_AGGREGATE_TOPIC", "articles-aggregate"),
		KafkaGroupID:     getEnv("KAFKA_GROUP_ID", "repo-article-group"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
