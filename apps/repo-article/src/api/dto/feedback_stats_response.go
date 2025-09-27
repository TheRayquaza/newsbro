package dto

type FeedbackStatsResponse struct {
	NewsID       uint    `json:"news_id"`
	ArticleTitle string  `json:"article_title,omitempty"`
	LikeCount    int64   `json:"like_count"`
	DislikeCount int64   `json:"dislike_count"`
	TotalCount   int64   `json:"total_count"`
	LikeRatio    float64 `json:"like_ratio"`
}
