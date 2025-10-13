package dto

type RSSUpdateRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Link        *string   `json:"link,omitempty"`
	Parents     *[]string `json:"parents,omitempty"`
}
