package aggregate

import (
	"time"
)

type ArticleAggregate struct {
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	Title       string    `json:"title"`
	Abstract    string    `json:"abstract"`
	Link        string    `json:"link"`
	Date        time.Time `json:"date"`
	IsActive    bool      `json:"is_active"`
}
