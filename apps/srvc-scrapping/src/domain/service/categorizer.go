package service

import (
	"regexp"
	"strings"

	"github.com/mmcdole/gofeed"
)

type ArticleCategorizer struct {
	categoryKeywords map[string][]string
	subCategoryMap   map[string]map[string][]string
}

func NewArticleCategorizer() *ArticleCategorizer {
	return &ArticleCategorizer{
		categoryKeywords: map[string][]string{
			"Technology":    {"tech", "software", "ai", "artificial intelligence", "computer", "digital", "app", "cloud", "cybersecurity", "blockchain"},
			"Business":      {"business", "economy", "market", "stock", "finance", "trade", "company", "startup", "investment"},
			"Politics":      {"politics", "government", "election", "congress", "parliament", "policy", "legislation", "democracy"},
			"Sports":        {"sports", "football", "basketball", "soccer", "tennis", "olympics", "athlete", "team", "championship"},
			"Health":        {"health", "medical", "disease", "hospital", "doctor", "medicine", "wellness", "fitness", "pandemic"},
			"Science":       {"science", "research", "study", "scientist", "discovery", "space", "astronomy", "biology", "physics"},
			"Entertainment": {"entertainment", "movie", "film", "music", "celebrity", "actor", "concert", "album", "show"},
			"Environment":   {"environment", "climate", "weather", "pollution", "sustainability", "renewable", "conservation"},
		},
		subCategoryMap: map[string]map[string][]string{
			"Technology": {
				"AI/ML":         {"artificial intelligence", "machine learning", "neural network", "deep learning", "ai"},
				"Cybersecurity": {"security", "hack", "breach", "vulnerability", "encryption", "malware"},
				"Software":      {"software", "app", "application", "code", "programming", "developer"},
				"Hardware":      {"hardware", "chip", "processor", "device", "smartphone", "computer"},
			},
			"Business": {
				"Finance":   {"finance", "bank", "investment", "stock", "trading", "wall street"},
				"Startups":  {"startup", "venture capital", "funding", "entrepreneur"},
				"Corporate": {"merger", "acquisition", "ceo", "earnings", "revenue"},
			},
			"Sports": {
				"Football":   {"football", "nfl", "quarterback", "super bowl"},
				"Basketball": {"basketball", "nba", "dunk", "playoff"},
				"Soccer":     {"soccer", "fifa", "premier league", "world cup"},
			},
			"Health": {
				"Medicine":      {"medicine", "drug", "treatment", "therapy", "clinical"},
				"Public Health": {"pandemic", "epidemic", "vaccination", "public health"},
				"Wellness":      {"wellness", "fitness", "nutrition", "mental health"},
			},
		},
	}
}

func (c *ArticleCategorizer) Categorize(item *gofeed.Item, description, content string) (string, string) {
	// Strategy 1: Use RSS feed categories
	if len(item.Categories) > 0 {
		category := normalizeCategory(item.Categories[0])
		subCategory := ""
		if len(item.Categories) > 1 {
			subCategory = normalizeCategory(item.Categories[1])
		}

		// If we found valid categories, try to map them
		mappedCategory := c.mapToStandardCategory(category)
		if mappedCategory != "" {
			return mappedCategory, subCategory
		}

		// If direct mapping failed, still use the original if no better match
		if category != "" {
			return category, subCategory
		}
	}

	// Strategy 2: Extract from URL path
	urlCategory, urlSubCategory := c.extractFromURL(item.Link)
	if urlCategory != "" {
		return urlCategory, urlSubCategory
	}

	// Strategy 3: Keyword matching in title, description, and content
	combinedText := strings.ToLower(item.Title + " " + description + " " + content)

	category := c.matchCategory(combinedText)
	subCategory := ""

	if category != "" {
		subCategory = c.matchSubCategory(category, combinedText)
	}

	// Strategy 4: Use feed metadata
	if category == "" && item.Custom != nil {
		if cat, ok := item.Custom["category"]; ok {
			category = normalizeCategory(cat)
		}
	}

	if category == "" {
		category = "General"
	}

	return category, subCategory
}

func (c *ArticleCategorizer) extractFromURL(url string) (string, string) {
	// Extract category from URL patterns like /category/subcategory/ or /section/subsection/
	re := regexp.MustCompile(`(?i)/(tech|business|politics|sports|health|science|entertainment|environment|world|local)/([^/]+)`)
	matches := re.FindStringSubmatch(url)

	if len(matches) > 1 {
		category := normalizeCategory(matches[1])
		subCategory := ""
		if len(matches) > 2 {
			subCategory = normalizeCategory(matches[2])
		}
		return c.mapToStandardCategory(category), subCategory
	}

	return "", ""
}

func (c *ArticleCategorizer) matchCategory(text string) string {
	bestMatch := ""
	maxScore := 0

	for category, keywords := range c.categoryKeywords {
		score := 0
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				score++
			}
		}
		if score > maxScore {
			maxScore = score
			bestMatch = category
		}
	}

	return bestMatch
}

func (c *ArticleCategorizer) matchSubCategory(category, text string) string {
	subCategories, exists := c.subCategoryMap[category]
	if !exists {
		return ""
	}

	bestMatch := ""
	maxScore := 0

	for subCat, keywords := range subCategories {
		score := 0
		for _, keyword := range keywords {
			if strings.Contains(text, keyword) {
				score++
			}
		}
		if score > maxScore {
			maxScore = score
			bestMatch = subCat
		}
	}

	return bestMatch
}

func (c *ArticleCategorizer) mapToStandardCategory(category string) string {
	category = strings.ToLower(category)

	mappings := map[string]string{
		"tech":          "Technology",
		"technology":    "Technology",
		"biz":           "Business",
		"business":      "Business",
		"finance":       "Business",
		"politics":      "Politics",
		"sport":         "Sports",
		"sports":        "Sports",
		"health":        "Health",
		"science":       "Science",
		"sci":           "Science",
		"entertainment": "Entertainment",
		"showbiz":       "Entertainment",
		"environment":   "Environment",
		"climate":       "Environment",
	}

	if mapped, ok := mappings[category]; ok {
		return mapped
	}

	return ""
}

func normalizeCategory(cat string) string {
	cat = strings.TrimSpace(cat)
	cat = strings.Title(strings.ToLower(cat))

	uuidPattern := regexp.MustCompile(`(?i)^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)
	if uuidPattern.MatchString(cat) {
		return ""
	}
	return cat
}
