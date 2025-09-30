package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"srvc_scrapping/src/config"
	"srvc_scrapping/src/data/models"
	"srvc_scrapping/src/data/repository"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/command"
	"github.com/mmcdole/gofeed"
)

type RSSService interface {
	ProcessFeed(ctx context.Context) (int, error)
}

type rSSService struct {
	articleRepo repository.ArticleRepository
	producer    sarama.SyncProducer
	config      *config.Config
	categorizer *ArticleCategorizer
}

func NewRSSService(articleRepo repository.ArticleRepository, producer sarama.SyncProducer, cfg *config.Config) RSSService {
	return &rSSService{
		articleRepo: articleRepo,
		producer:    producer,
		config:      cfg,
		categorizer: NewArticleCategorizer(),
	}
}

func (u *rSSService) ProcessFeed(ctx context.Context) (int, error) {
	fp := gofeed.NewParser()
	totalCount := 0
	processedCount := 0

	for _, feedURL := range u.config.RSSFeedURL {
		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			log.Printf("Error parsing feed %s: %v", feedURL, err)
			continue
		}

		for _, item := range feed.Items {
			exists, err := u.articleRepo.Exists(ctx, item.Link)
			if err != nil {
				log.Printf("Error checking if article exists: %v", err)
				continue
			}
			if exists {
				continue
			}

			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}

			description := cleanHTML(item.Description)
			content := cleanHTML(item.Content)

			// Extract categories using multiple strategies
			category, subCategory := u.categorizer.Categorize(item, description, content)

			article := &models.Article{
				Link: item.Link,
			}

			if err := u.articleRepo.Create(ctx, article); err != nil {
				log.Printf("Error saving article: %v", err)
				continue
			}

			message := command.NewArticleCommand{
				Title:       item.Title,
				Link:        item.Link,
				Description: description,
				Content:     content,
				Author:      getAuthor(item),
				PublishedAt: publishedAt,
				Category:    category,
				Subcategory: subCategory,
			}

			if err := u.sendToKafka(&message); err != nil {
				log.Printf("Error sending to Kafka: %v", err)
				continue
			}

			processedCount++
			totalCount++
		}

		if processedCount > 0 {
			log.Printf("Processed %d new articles from feed %s", processedCount, feedURL)
			processedCount = 0
		}
	}

	return totalCount, nil
}

// ArticleCategorizer handles category extraction
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

func (u *rSSService) sendToKafka(message *command.NewArticleCommand) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(message.Link),
		Topic: u.config.KafkaTopic,
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = u.producer.SendMessage(msg)
	return err
}

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

	cleaned = strings.ReplaceAll(cleaned, """, "\"")
	cleaned = strings.ReplaceAll(cleaned, """, "\"")
	cleaned = strings.ReplaceAll(cleaned, "'", "'")
	cleaned = strings.ReplaceAll(cleaned, "'", "'")

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
