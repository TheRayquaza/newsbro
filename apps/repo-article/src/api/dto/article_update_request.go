package dto

type ArticleUpdateRequest struct {
	ID          uint    `json:"id" binding:"required"`
	Category    *string `json:"category,omitempty"`
	Subcategory *string `json:"subcategory,omitempty"`
	Title       *string `json:"title,omitempty"`
	Abstract    *string `json:"abstract,omitempty"`
	Link        *string `json:"link,omitempty"`
}
