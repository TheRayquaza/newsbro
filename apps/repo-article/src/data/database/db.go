package database

import (
	"repo_article/src/data/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/TheRayquaza/newsbro/apps/libs/utils"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		utils.SugarLog.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	utils.SugarLog.Info("Connected to database successfully")

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.Article{},
		&models.Feedback{},
		&models.RSSSource{},
	)
	if err != nil {
		utils.SugarLog.Errorf("Failed to auto-migrate database: %v", err)
		return nil, err
	}

	utils.SugarLog.Info("Database auto-migration completed successfully")

	return db, nil
}
