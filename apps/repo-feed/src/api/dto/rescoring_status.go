package dto

import "time"

type RescoringStatus struct {
	CurrentCursor uint64        `json:"current_cursor"`
	BatchSIze     int           `json:"batch_size"`
	Interval      time.Duration `json:"interval"`
	Workers       int           `json:"workers"`
	Threshold     float64       `json:"threshold"`
	Enabled       bool          `json:"enabled"`
}
