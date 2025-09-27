package config

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DatabaseURL string
	JWTSecret   string
	Admins	  	[]string
}


func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	admins := strings.Split(getEnv("ADMINS", ""), ",")
	filtered := make([]string, 0, len(admins))
	for _, a := range admins {
		if trimmed := strings.TrimSpace(a); trimmed != "" {
			filtered = append(filtered, trimmed)
		}
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
		DatabaseURL: fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
			getEnv("DATABASE_USERNAME", "username"),
			getEnv("DATABASE_PASSWORD", "password"),
			getEnv("DATABASE_HOST", "localhost:5432"),
			getEnv("DATABASE_NAME", "repo_account"),
		),
		JWTSecret: getEnv("JWT_SECRET", ""),
		Admins:    filtered,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
