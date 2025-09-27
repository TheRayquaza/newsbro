package entities

import "time"

type FeedbackCSVExport struct {
	FeedbackID uint      `json:"feedback_id"`
	UserID     uint      `json:"user_id"`
	NewsID     uint      `json:"news_id"`
	Value      int       `json:"value"`
	ValueText  string    `json:"value_text"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
