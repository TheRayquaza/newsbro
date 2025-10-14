package dto

import "time"

type TreeRSSResponse struct {
	Name        string            `json:"name"`
	DisplayName string            `json:"display_name"`
	Description string            `json:"description"`
	Link        string            `json:"link"`
	Active      bool              `json:"active"`
	Children    []TreeRSSResponse `json:"children,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}
