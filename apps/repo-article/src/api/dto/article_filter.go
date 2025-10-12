package dto

import "time"

type ArticleFilters struct {
	Category    string     `form:"category"`
	Subcategory string     `form:"subcategory"`
	Search      string     `form:"search"`
	BeginDate   *time.Time `form:"begin_date" time_format:"2006-01-02" time_utc:"1"`
	EndDate     *time.Time `form:"end_date" time_format:"2006-01-02" time_utc:"1"`
	Limit       uint       `form:"limit,default=10"`
	Offset      uint       `form:"offset,default=0"`
}
