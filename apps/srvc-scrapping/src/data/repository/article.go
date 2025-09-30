package repository

import (
	"context"
	"gorm.io/gorm"
	"srvc_scrapping/src/data/models"
)

type ArticleRepository interface {
	Exists(ctx context.Context, link string) (bool, error)
	Create(ctx context.Context, article *models.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) Exists(ctx context.Context, link string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Article{}).Where("link = ?", link).Count(&count).Error
	return count > 0, err
}

func (r *articleRepository) Create(ctx context.Context, article *models.Article) error {
	dbArticle := &models.Article{
		Link: article.Link,
	}
	return r.db.WithContext(ctx).Create(dbArticle).Error
}
