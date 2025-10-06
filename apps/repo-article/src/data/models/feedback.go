package models

import (
	"time"
)

type Feedback struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"not null;index:idx_user_article,unique" json:"user_id"`
	NewsID    uint      `gorm:"not null;index:idx_user_article,unique;index:idx_news_feedback" json:"news_id"`
	Article   Article   `gorm:"foreignKey:NewsID;references:ID" json:"article"` // <-- link to Article
	Value     int       `gorm:"not null;check:value IN (0, 1)" json:"value"`    // 0 = dislike, 1 = like
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Feedback) TableName() string {
	return "feedbacks"
}
