package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port             string
	JWTSecret        string
	FrontendOrigin   string
	LoginRedirectURL string
	// Redis
	RedisSentinels  []string
	RedisMasterName string
	RedisPassword   string
	RedisDB         int
	// Kafka
	KafkaBrokers               []string
	KafkaInferenceCommandTopic string
	KafkaGroupID               string
	// Other
	DefaultModel string
	Environment  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	redisSentinels := func() []string {
		sentinelsStr := getEnv("REDIS_SENTINELS", "localhost:26379")
		return splitAndTrim(sentinelsStr, ",")
	}()

	kafkaBrokers := func() []string {
		brokersStr := getEnv("KAFKA_BROKERS", "localhost:9092")
		return splitAndTrim(brokersStr, ",")
	}()

	redisDb := func() int {
		dbStr := getEnv("REDIS_DB", "0")
		db, err := strconv.Atoi(dbStr)
		if err != nil {
			fmt.Printf("Invalid REDIS_DB value '%s', defaulting to 0\n", dbStr)
			return 0
		}
		return db
	}()

	return &Config{
		Port:                       getEnv("PORT", "8080"),
		JWTSecret:                  getEnv("JWT_SECRET", "your-secret-key"),
		FrontendOrigin:             getEnv("FRONTEND_ORIGIN", "http://localhost:3000"),
		LoginRedirectURL:           getEnv("LOGIN_REDIRECT_URL", "http://localhost:3000/auth/callback"),
		RedisSentinels:             redisSentinels,
		RedisMasterName:            getEnv("REDIS_MASTER_NAME", "mymaster"),
		RedisPassword:              getEnv("REDIS_PASSWORD", ""),
		RedisDB:                    redisDb,
		KafkaBrokers:               kafkaBrokers,
		KafkaInferenceCommandTopic: getEnv("KAFKA_INFERENCE_COMMAND_TOPIC", "inference-commands"),
		KafkaGroupID:               getEnv("KAFKA_GROUP_ID", "repo-feed-group"),
		DefaultModel:               getEnv("DEFAULT_MODEL", "tfidf"),
		Environment:                getEnv("ENVIRONMENT", "prod"),
	}
}

func splitAndTrim(s, sep string) []string {
	parts := strings.Split(s, sep)
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
