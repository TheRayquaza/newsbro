package models

import (
	"time"
)

type ArticleModel struct {
	ID          uint      `json:"id"`
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	Title       string    `json:"title"`
	Abstract    string    `json:"abstract"`
	Link        string    `json:"link"`
	RSSLink     string    `json:"rss_link"`
	PublishedAt time.Time `json:"published_at"`
	Score       float64   `json:"score"`
	DecayScore  float64   `json:"decay_score"`
}
