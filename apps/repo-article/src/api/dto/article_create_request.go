package dto

type ArticleCreateRequest struct {
	Category    string `json:"category" binding:"required"`
	Subcategory string `json:"subcategory" binding:"required"`
	Title       string `json:"title" binding:"required"`
	Abstract    string `json:"abstract"`
	Link        string `json:"link" binding:"required,url"`
}
