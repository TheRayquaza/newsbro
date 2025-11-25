package main

import (
	"repo_account/src/api/routes"
	"repo_account/src/config"
	"repo_account/src/data/database"
	"repo_account/src/domain/services"

	"github.com/TheRayquaza/newsbro/apps/libs/utils"
)

func main() {
	cfg := config.Load()

	if err := utils.Initialize(cfg.Environment); err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}
	defer utils.SugarLog.Sync()

	if cfg.Environment == "dev" {
		utils.SugarLog.Info("Running in development mode")
	}

	utils.SugarLog.Infof("Connecting to database at %s", cfg.DatabaseURL)
	db, err := database.Initialize(cfg.DatabaseURL)
	if err != nil {
		utils.SugarLog.Fatal("Failed to initialize database:", err)
	}

	utils.SugarLog.Debug("Initializing services")
	authService := services.NewAuthService(cfg, db)
	userService := services.NewUserService(db)

	utils.SugarLog.Debug("Setting up routes")
	router := routes.SetupRouter(cfg, authService, userService)

	utils.SugarLog.Infof("Server starting on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		utils.SugarLog.Fatal("Failed to start server:", err)
	}
}
