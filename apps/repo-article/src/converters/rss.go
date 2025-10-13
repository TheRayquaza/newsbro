package converters

import (
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
)

func RSSModelToRSSResponse(rss *models.RSSSource) dto.RSSResponse {
	return dto.RSSResponse{
		Name:        rss.Name,
		Description: rss.Description,
		Link:        rss.Link,
		Active:      rss.Active,
		CreatedAt:   rss.CreatedAt,
		UpdatedAt:   rss.UpdatedAt,
	}
}

func RSSModelToRSSAggregate(rss *models.RSSSource, active bool) aggregate.RSSAggregate {
	return aggregate.RSSAggregate{
		Name:        rss.Name,
		Description: rss.Description,
		Link:        rss.Link,
		Active:      active,
	}
}
