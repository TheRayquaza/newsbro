package main

import (
	"log"
	"repo_account/src/api/routes"
	"repo_account/src/config"
	"repo_account/src/data/database"
	"repo_account/src/domain/services"
)

func main() {
	// Load configuration
	cfg := config.Load()

	if cfg.Environment == "dev" {
		log.Println("Running in development mode")
	}

	// Initialize database
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize services
	authService := services.NewAuthService(cfg, db)
	userService := services.NewUserService(db)

	// Setup routes
	router := routes.SetupRouter(cfg, authService, userService)

	// Start server
	log.Printf("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
