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
	RSSLink     string    `json:"rss_link"`
	PublishedAt time.Time `json:"published_at"`
	IsActive    bool      `json:"is_active"`
}
