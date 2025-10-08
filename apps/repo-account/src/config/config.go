package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	Port                    string
	DatabaseURL             string
	JWTSecret               string
	OIDCIssuerURL           string
	OIDCClientID            string
	OIDCClientSecret        string
	OIDCRedirectURL         string
	OIDCRedirectFrontendURL string
	CookieDomain            string
	LoginRedirectURL        string
	FrontendDomain          string
	Environment             string
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
		JWTSecret:               getEnv("JWT_SECRET", "your-secret-key"),
		OIDCIssuerURL:           getEnv("OIDC_ISSUER_URL", ""),
		OIDCClientID:            getEnv("OIDC_CLIENT_ID", ""),
		OIDCClientSecret:        getEnv("OIDC_CLIENT_SECRET", ""),
		OIDCRedirectURL:         getEnv("OIDC_REDIRECT_URL", "http://localhost:8080/auth/callback"),
		CookieDomain:            getEnv("COOKIE_DOMAIN", ""),
		OIDCRedirectFrontendURL: getEnv("OIDC_REDIRECT_FRONTEND_URL", "http://localhost:3000"),
		LoginRedirectURL:        getEnv("LOGIN_REDIRECT_URL", "http://localhost:3000"),
		FrontendOrigin:          getEnv("FRONTEND_ORIGIN", "https://app.newsbro.cc"),
		Environment:             getEnv("ENVIRONMENT", "prod"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
