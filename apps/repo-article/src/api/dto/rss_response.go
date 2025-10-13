package dto

import "time"

type RSSResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Link        string    `json:"link"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
