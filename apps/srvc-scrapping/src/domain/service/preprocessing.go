package service

import (
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

func cleanHTML(content string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	cleaned := re.ReplaceAllString(content, "")

	cleaned = strings.ReplaceAll(cleaned, "&nbsp;", " ")
	cleaned = strings.ReplaceAll(cleaned, "&amp;", "&")
	cleaned = strings.ReplaceAll(cleaned, "&lt;", "<")
	cleaned = strings.ReplaceAll(cleaned, "&gt;", ">")
	cleaned = strings.ReplaceAll(cleaned, "&quot;", "\"")
	cleaned = strings.ReplaceAll(cleaned, "&#39;", "'")

	re = regexp.MustCompile(`([.!?,;:]){2,}`)
	cleaned = re.ReplaceAllString(cleaned, "$1")

	re = regexp.MustCompile(`\s[•◦▪▫■□●○★☆→←↑↓]+\s`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	re = regexp.MustCompile(`[-_]{3,}`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	re = regexp.MustCompile(`^[.,;:!?]+|[.,;:!?]+$`)
	cleaned = re.ReplaceAllString(cleaned, "")

	re = regexp.MustCompile(`\s+[.,;:!?]+\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	re = regexp.MustCompile(`\.{2,}|…`)
	cleaned = re.ReplaceAllString(cleaned, ".")

	re = regexp.MustCompile(`\s+([.,;:!?])`)
	cleaned = re.ReplaceAllString(cleaned, "$1")

	cleaned = strings.TrimSpace(cleaned)
	re = regexp.MustCompile(`\s+`)
	cleaned = re.ReplaceAllString(cleaned, " ")

	return cleaned
}

func getAuthor(item *gofeed.Item) string {
	if item.Author != nil {
		return item.Author.Name
	}
	for _, author := range item.Authors {
		if author.Name != "" {
			return author.Name
		}
	}
	return ""
}
