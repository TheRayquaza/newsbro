package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
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
	KafkaFeedbackTopic         string
	KafkaGroupID               string
	KafkaFeedbackGroupID       string

	// Feedback settings
	DefaultModel       string
	Models             []string
	FeedbackExpiration time.Duration

	// Feed Rescoring
	FeedRescoring struct {
		Enabled           bool
		Interval          time.Duration
		BatchSize         int
		ConcurrentWorkers int
		ScoreThreshold    float64
		DecayEnabled      bool
		DecayHalfLife     time.Duration
	}

	// Environment
	Environment string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// --- Redis ---
	redisSentinels := splitAndTrim(getEnv("REDIS_SENTINELS", "localhost:26379"), ",")
	redisDB := parseIntEnv("REDIS_DB", 0)

	// --- Kafka ---
	kafkaBrokers := splitAndTrim(getEnv("KAFKA_BROKERS", "localhost:9092"), ",")

	// --- Feedback expiration ---
	feedbackExpiration := func() time.Duration {
		expStr := getEnv("FEEDBACK_EXPIRATION_SECONDS", "604800") // Default 7 days
		seconds, err := strconv.Atoi(expStr)
		if err != nil {
			fmt.Printf("Invalid FEEDBACK_EXPIRATION_SECONDS='%s', using 604800\n", expStr)
			seconds = 604800
		}
		return time.Duration(seconds) * time.Second
	}()

	// --- Feed Rescoring ---
	feedRescoring := struct {
		Enabled           bool
		Interval          time.Duration
		BatchSize         int
		ConcurrentWorkers int
		ScoreThreshold    float64
		DecayEnabled      bool
		DecayHalfLife     time.Duration
	}{
		Enabled:           parseBoolEnv("FEED_RESCORING_ENABLED", false),
		Interval:          parseDurationEnv("FEED_RESCORING_INTERVAL", "10m"),
		BatchSize:         parseIntEnv("FEED_RESCORING_BATCH_SIZE", 1000),
		ConcurrentWorkers: parseIntEnv("FEED_RESCORING_WORKERS", 10),
		ScoreThreshold:    parseFloatEnv("FEED_RESCORING_THRESHOLD", 0.1),
		DecayEnabled:      parseBoolEnv("FEED_RESCORING_DECAY_ENABLED", true),
		DecayHalfLife:     parseDurationEnv("FEED_RESCORING_DECAY_HALFLIFE", "120h"), // default 5 days
	}

	// --- Return Config ---
	return &Config{
		Port:             getEnv("PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "your-secret-key"),
		FrontendOrigin:   getEnv("FRONTEND_ORIGIN", "https://app.newsbro.cc"),
		LoginRedirectURL: getEnv("LOGIN_REDIRECT_URL", "http://localhost:3000/auth/callback"),

		RedisSentinels:  redisSentinels,
		RedisMasterName: getEnv("REDIS_MASTER_NAME", "mymaster"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         redisDB,

		KafkaBrokers:               kafkaBrokers,
		KafkaInferenceCommandTopic: getEnv("KAFKA_INFERENCE_COMMAND_TOPIC", "inference-commands"),
		KafkaFeedbackTopic:         getEnv("KAFKA_FEEDBACK_AGGREGATE_TOPIC", "feedback-aggregate"),
		KafkaGroupID:               getEnv("KAFKA_GROUP_ID", "repo-feed-group"),
		KafkaFeedbackGroupID:       getEnv("KAFKA_FEEDBACK_GROUP_ID", "repo-feed-feedback-group"),

		DefaultModel:       getEnv("DEFAULT_MODEL", "tfidf"),
		Models:             splitAndTrim(getEnv("AVAILABLE_MODELS", "tfidf,sbert"), ","),
		FeedbackExpiration: feedbackExpiration,

		FeedRescoring: feedRescoring,

		Environment: getEnv("ENVIRONMENT", "prod"),
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

func parseBoolEnv(key string, defaultValue bool) bool {
	val := strings.ToLower(getEnv(key, fmt.Sprintf("%t", defaultValue)))
	return val == "true" || val == "1" || val == "yes"
}

func parseIntEnv(key string, defaultValue int) int {
	valStr := getEnv(key, fmt.Sprintf("%d", defaultValue))
	val, err := strconv.Atoi(valStr)
	if err != nil {
		log.Printf("Invalid %s='%s', using default %d\n", key, valStr, defaultValue)
		return defaultValue
	}
	return val
}

func parseFloatEnv(key string, defaultValue float64) float64 {
	valStr := getEnv(key, fmt.Sprintf("%f", defaultValue))
	val, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		log.Printf("Invalid %s='%s', using default %f\n", key, valStr, defaultValue)
		return defaultValue
	}
	return val
}

func parseDurationEnv(key, defaultValue string) time.Duration {
	valStr := getEnv(key, defaultValue)
	dur, err := time.ParseDuration(valStr)
	if err != nil {
		log.Printf("Invalid %s='%s', using default %s\n", key, valStr, defaultValue)
		dur, _ = time.ParseDuration(defaultValue)
	}
	return dur
}
