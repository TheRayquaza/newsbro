package aggregate

type FeedbackAggregate struct {
	UserID   uint `json:"user_id"`
	NewsID   uint `json:"news_id"`
	Value    int  `json:"value"`     // 0 = dislike, 1 = like
	IsActive bool `json:"is_active"` // remove feedback if false
}
