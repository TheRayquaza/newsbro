package repository

import (
	"context"
	"gorm.io/gorm"
	"srvc_scrapping/src/data/models"
)

type RSSRepository interface {
	Exists(ctx context.Context, link string) (bool, error)
	Create(ctx context.Context, rss *models.RSS) error
	GetAllLinks(ctx context.Context) ([]string, error)
	DeleteByLink(ctx context.Context, link string) error
}

type rssRepository struct {
	db *gorm.DB
}

func NewRSSRepository(db *gorm.DB) RSSRepository {
	return &rssRepository{db: db}
}

func (r *rssRepository) Exists(ctx context.Context, link string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.RSS{}).Where("link = ?", link).Count(&count).Error
	return count > 0, err
}

func (r *rssRepository) Create(ctx context.Context, rss *models.RSS) error {
	dbRSS := &models.RSS{
		Link: rss.Link,
	}
	return r.db.WithContext(ctx).Create(dbRSS).Error
}

func (r *rssRepository) GetAllLinks(ctx context.Context) ([]string, error) {
	var rsss []models.RSS
	err := r.db.WithContext(ctx).Find(&rsss).Error
	links := make([]string, len(rsss))
	for i, rss := range rsss {
		links[i] = rss.Link
	}
	return links, err
}

func (r *rssRepository) DeleteByLink(ctx context.Context, link string) error {
	return r.db.WithContext(ctx).Where("link = ?", link).Delete(&models.RSS{}).Error
}
