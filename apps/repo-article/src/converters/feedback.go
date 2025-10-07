package converters

import (
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"repo_article/src/api/dto"
	"repo_article/src/data/models"
)

func FeedbackModelToFeedbackResponse(feedback *models.Feedback) dto.FeedbackResponse {
	return dto.FeedbackResponse{
		UserID: feedback.UserID,
		NewsID: feedback.NewsID,
		Value:  feedback.Value,
		Date:   feedback.UpdatedAt,
	}
}

func FeedbackModelToFeedbackAggregate(feedback *models.Feedback, active bool) aggregate.FeedbackAggregate {
	return aggregate.FeedbackAggregate{
		UserID:   feedback.UserID,
		NewsID:   feedback.NewsID,
		Value:    feedback.Value,
		Date:     feedback.UpdatedAt,
		IsActive: active,
	}
}
