package converters

import (
	"repo_feed/src/api/dto"
	"repo_feed/src/data/models"
)

func FeedModelToFeedResponse(feed *models.FeedModel) dto.FeedResponse {
	return dto.FeedResponse{
		UserID:    feed.UserID,
		ModelName: feed.ModelName,
		Articles: func() []dto.Article {
			articles := make([]dto.Article, len(feed.Articles))
			for i, article := range feed.Articles {
				articles[i] = ArticleModelToArticle(&article)
			}
			return articles
		}(),
		UpdatedAt: feed.UpdatedAt,
	}
}
