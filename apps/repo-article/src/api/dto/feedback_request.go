package dto

type FeedbackRequest struct {
	Value int `json:"value" binding:"required,min=0,max=1" example:"1"`
}
