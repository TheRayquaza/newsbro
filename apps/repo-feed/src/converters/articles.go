package converters

import (
	"repo_feed/src/api/dto"
	"repo_feed/src/data/models"
)

func ArticleModelToArticle(article *models.ArticleModel) dto.Article {
	return dto.Article{
		ID:          article.ID,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Link:        article.Link,
		RSSLink:     article.RSSLink,
		PublishedAt: article.PublishedAt,
		Score:       article.DecayScore,
	}
}

func ArticleToArticleModel(article *dto.Article, score float64) models.ArticleModel {
	return models.ArticleModel{
		ID:          article.ID,
		Category:    article.Category,
		Subcategory: article.Subcategory,
		Title:       article.Title,
		Abstract:    article.Abstract,
		Link:        article.Link,
		RSSLink:     article.RSSLink,
		PublishedAt: article.PublishedAt,
		Score:       score,
		DecayScore:  score,
	}
}
