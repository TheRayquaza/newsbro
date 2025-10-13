package dto

type RSSCreateRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Link        string    `json:"link" binding:"omitempty,url"`
	Parents     *[]string `json:"parents,omitempty"`
}
