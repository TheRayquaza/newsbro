package dto

import "time"

type RSSFilters struct {
	BeginDate *time.Time `form:"begin_date" time_format:"2006-01-02" time_utc:"1"`
	EndDate   *time.Time `form:"end_date" time_format:"2006-01-02" time_utc:"1"`
}
