package dto

import "time"

type ArticleFilters struct {
	Category    string     `form:"category"`
	Subcategory string     `form:"subcategory"`
	Search      string     `form:"search"`
	FeedName    string     `form:"feed_name"`
	BeginDate   *time.Time `form:"begin_date" time_format:"2006-01-02" time_utc:"1"`
	EndDate     *time.Time `form:"end_date" time_format:"2006-01-02" time_utc:"1"`
	Count       bool       `form:"count,default=false"`
	Limit       uint       `form:"limit,default=10"`
	Offset      uint       `form:"offset,default=0"`
	SortBy      string     `form:"sort_by,default=none"`
}
