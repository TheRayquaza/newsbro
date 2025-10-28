package converters

import (
	"repo_feed/src/api/dto"
	"repo_feed/src/data/models"
)

func ArticleModelToArticleResponse(article *models.ArticleModel) dto.Article {
	return dto.Article{
		ID:          article.ID,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Link:        article.Link,
		RSSLink:     article.RSSLink,
		PublishedAt: article.PublishedAt,
		IsActive:    article.IsActive,
	}
}

func ArticleToArticleModel(article *dto.Article) models.ArticleModel {
	return models.ArticleModel{
		ID:          article.ID,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Link:        article.Link,
		RSSLink:     article.RSSLink,
		PublishedAt: article.PublishedAt,
		IsActive:    true, // Default to active when creating from DTO
	}
}
