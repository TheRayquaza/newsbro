package dto

import "time"

type FeedProcessingStats struct {
	FeedURL         string
	TotalItems      int
	ProcessedCount  int
	SkippedCount    int
	ErrorCount      int
	ErrorDetails    []string
	ProcessingTime  time.Duration
	LanguageSkipped int
}
