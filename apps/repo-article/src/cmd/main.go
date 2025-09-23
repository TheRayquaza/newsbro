package main

import (
	"log"
	"repo_article/src/api/routes"
	"repo_article/src/config"
	"repo_article/src/data/database"
	"repo_article/src/domain/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	articleService := services.NewArticleService(db)
	authService := services.NewAuthService()

	// Setup routes
	router := routes.SetupRouter(cfg, articleService, authService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
