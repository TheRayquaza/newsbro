package dto

import "time"

type ArticleResponse struct {
	ID              uint      `json:"id"`
	Category        string    `json:"category"`
	Subcategory     string    `json:"subcategory"`
	Title           string    `json:"title"`
	Abstract        string    `json:"abstract"`
	Link            string    `json:"link" gorm:"not null;uniqueIndex"`
	Value           int       `json:"liked_value"`
	LastInteraction time.Time `json:"last_interaction"`
	PublishedAt     time.Time `json:"published_at"`
}
