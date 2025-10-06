package database

import (
	"repo_article/src/data/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Initialize(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto-migrate models
	err = db.AutoMigrate(
		&models.Article{},
		&models.Feedback{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}
