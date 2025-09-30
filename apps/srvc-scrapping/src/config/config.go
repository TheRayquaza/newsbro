package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	KafkaBrokers []string
	KafkaTopic   string
	RSSFeedURL   []string
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
		Port: getEnv("PORT", "8080"),
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			getEnv("DATABASE_USERNAME", "username"),
			getEnv("DATABASE_PASSWORD", "password"),
			getEnv("DATABASE_HOST", "localhost:5432"),
			getEnv("DATABASE_NAME", "repo_account"),
		),
		KafkaBrokers: strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
		KafkaTopic:   getEnv("KAFKA_TOPIC", "new-articles-command"),
		RSSFeedURL:   Map(strings.FieldsFunc(getEnv("RSS_FEED_URLS", ""), func(r rune) bool { return r == ',' || r == '\n' }), func(s string) string { return strings.TrimSpace(s) }),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
