package command

import (
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"time"
)

type InferenceCommand struct {
	UserID  uint      `json:"user_id"`
	Model   string    `json:"model"`
	Score   float64   `json:"score"`
	Date    time.Time `json:"date"`
	Article *aggregate.ArticleAggregate `json:"article"`
}
