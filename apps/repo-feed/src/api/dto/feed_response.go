package dto

import (
	"time"
)

type FeedResponse struct {
	UserID    uint      `json:"user_id"`
	ModelName string    `json:"model_name"`
	Articles  []Article `json:"articles"`
	UpdatedAt time.Time `json:"updated_at"`
}
