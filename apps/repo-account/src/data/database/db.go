package database

import (
	"repo_account/src/data/models"

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
	err = db.AutoMigrate(
		&models.User{},
		&models.RefreshToken{},
	)
	if err != nil {
		utils.SugarLog.Errorf("Failed to auto-migrate database: %v", err)
		return nil, err
	}
	utils.SugarLog.Info("Database auto-migration completed successfully")

	return db, nil
}
