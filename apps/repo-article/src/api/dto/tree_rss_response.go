package dto

import "time"

type TreeRSSResponse struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	Active      bool              `json:"active"`
	Children    []TreeRSSResponse `json:"children,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
