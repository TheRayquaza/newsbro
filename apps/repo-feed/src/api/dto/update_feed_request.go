package dto

import (
	"time"
)

type UpdateFeedRequest struct {
	UserID  uint      `json:"user_id"`
	Model   string    `json:"model_name"`
	Score   float64   `json:"score"`
	Article *Article  `json:"article"`
	Date    time.Time `json:"date"`
}
