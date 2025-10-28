package models

import (
	"time"
)

type FeedModel struct {
	UserID    uint           `json:"user_id"`
	ModelName string         `json:"model_name"`
	Articles  []ArticleModel `json:"articles"`
	UpdatedAt time.Time      `json:"updated_at"`
}
