package aggregate

import (
	"time"
)

type ArticleAggregate struct {
	ID          uint      `json:"id"`
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	Title       string    `json:"title"`
	Abstract    string    `json:"abstract"`
	Link        string    `json:"link"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	IsActive    bool      `json:"is_active"`
}
