package command

import "time"

type NewArticleCommand struct {
	Title       string    `json:"title"`
	Link        string    `json:"link"`
	Description string    `json:"description"`
	Content     string    `json:"content"`
	Author      string    `json:"author"`
	Category    string    `json:"category"`
	Subcategory string    `json:"subcategory"`
	PublishedAt time.Time `json:"published_at"`
}
