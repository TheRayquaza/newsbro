package dto

import "time"

type ArticleCreateRequest struct {
	ID          uint      `json:"id,omitempty"` // Optional: If provided, must be unique
	Category    string    `json:"category" binding:"required"`
	Subcategory string    `json:"subcategory" binding:"required"`
	Title       string    `json:"title" binding:"required"`
	Abstract    string    `json:"abstract"`
	Link        string    `json:"link" binding:"required,url"`
	PublishedAt time.Time `json:"published_at" binding:"required"`
}
