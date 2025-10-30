package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

const zMemberDelimiter = "|"

func encodeZSetMember(articleID uint, score float64, publishedAt time.Time) string {
	return fmt.Sprintf("%d%s%.6f%s%d", articleID, zMemberDelimiter, score, zMemberDelimiter, publishedAt.Unix())
}

func decodeZSetMember(member string) (uint, float64, time.Time, error) {
	parts := strings.Split(member, zMemberDelimiter)
	if len(parts) != 3 {
		return 0, 0, time.Time{}, fmt.Errorf("invalid ZSET member format: expected 3 parts, got %d", len(parts))
	}

	articleIDUint64, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to parse article ID '%s': %w", parts[0], err)
	}
	articleID := uint(articleIDUint64)

	originalScore, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to parse original score '%s': %w", parts[1], err)
	}

	publishedAtUnix, err := strconv.ParseInt(parts[2], 10, 64)
	if err != nil {
		return 0, 0, time.Time{}, fmt.Errorf("failed to parse publishedAt Unix time '%s': %w", parts[2], err)
	}
	publishedAt := time.Unix(publishedAtUnix, 0)

	return articleID, originalScore, publishedAt, nil
}
