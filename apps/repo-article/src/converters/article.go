package converters

import (
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
)

func ArticleModelToArticleResponse(article *models.Article) dto.ArticleResponse {
	return dto.ArticleResponse{
		ID:          article.ID,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Link:        article.Link,
		RSSLink:     article.RSS.Link,
		PublishedAt: article.PublishedAt,
	}
}

func ArticleModelToArticleAggregate(article *models.Article, active bool) aggregate.ArticleAggregate {
	return aggregate.ArticleAggregate{
		ID:          article.ID,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Link:        article.Link,
		RSSLink:     article.RSS.Link,
		PublishedAt: article.PublishedAt,
		IsActive:    active,
	}
}
