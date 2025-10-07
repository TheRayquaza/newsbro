package dto

import "time"

type FeedbackResponse struct {
	ID     uint      `json:"id"`
	UserID uint      `json:"user_id"`
	NewsID uint      `json:"news_id"`
	Value  int       `json:"value"`
	Date   time.Time `json:"date"`
}
