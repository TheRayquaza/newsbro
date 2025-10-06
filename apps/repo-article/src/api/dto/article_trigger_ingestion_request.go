package dto

import "time"

type ArticleTriggerIngestionRequest struct {
	BeginDate time.Time `json:"begin_date" binding:"required"`
	EndDate   time.Time `json:"end_date" binding:"required"`
}
