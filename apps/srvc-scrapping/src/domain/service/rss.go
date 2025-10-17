package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"srvc_scrapping/src/api/dto"
	"srvc_scrapping/src/config"
	"srvc_scrapping/src/data/models"
	"srvc_scrapping/src/data/repository"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/command"
	"github.com/gtuk/discordwebhook"
	"github.com/mmcdole/gofeed"
	"github.com/pemistahl/lingua-go"
)

type RSSService struct {
	articleRepo repository.ArticleRepository
	rssRepo     repository.RSSRepository
	producer    sarama.SyncProducer
	cfg         *config.Config
	categorizer *ArticleCategorizer
	detector    lingua.LanguageDetector
}

func NewRSSService(articleRepo repository.ArticleRepository, rssRepo repository.RSSRepository, producer sarama.SyncProducer, cfg *config.Config) *RSSService {
	languages := []lingua.Language{
		lingua.Arabic, lingua.Chinese, lingua.English,
		lingua.French, lingua.German, lingua.Hindi, lingua.Japanese,
		lingua.Portuguese, lingua.Russian, lingua.Spanish, lingua.Swedish,
	}

	detector := lingua.NewLanguageDetectorBuilder().
		FromLanguages(languages...).
		Build()

	return &RSSService{
		articleRepo: articleRepo,
		rssRepo:     rssRepo,
		producer:    producer,
		cfg:         cfg,
		categorizer: NewArticleCategorizer(),
		detector:    detector,
	}
}

func (u *RSSService) ProcessFeed(ctx context.Context) (int, error) {
	fp := gofeed.NewParser()
	startTime := time.Now()

	totalCount := 0
	feedStats := make([]dto.FeedProcessingStats, 0)
	globalErrors := make([]string, 0)

	RSSFeedURL, err := u.rssRepo.GetAllLinks(ctx)
	if err != nil {
		log.Printf("Error retrieving RSS feed URLs: %v", err)
		return 0, fmt.Errorf("failed to retrieve RSS feed URLs: %w", err)
	}

	for _, feedURL := range RSSFeedURL {
		feedStartTime := time.Now()
		stats := dto.FeedProcessingStats{
			FeedURL:      feedURL,
			ErrorDetails: make([]string, 0),
		}

		feed, err := fp.ParseURL(feedURL)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to parse feed: %v", err)
			stats.ErrorDetails = append(stats.ErrorDetails, errMsg)
			globalErrors = append(globalErrors, fmt.Sprintf("[%s] %s", feedURL, errMsg))
			log.Printf("Error parsing feed %s: %v", feedURL, err)
			feedStats = append(feedStats, stats)
			continue
		}

		stats.TotalItems = len(feed.Items)

		for _, item := range feed.Items {
			exists, err := u.articleRepo.Exists(ctx, item.Link)
			if err != nil {
				errMsg := fmt.Sprintf("DB check failed for %s: %v", item.Link, err)
				stats.ErrorDetails = append(stats.ErrorDetails, errMsg)
				stats.ErrorCount++
				log.Printf("Error checking if article exists: %v", err)
				continue
			}

			if exists {
				stats.SkippedCount++
				continue
			}

			publishedAt := time.Now()
			if item.PublishedParsed != nil {
				publishedAt = *item.PublishedParsed
			}

			description := cleanHTML(item.Description)
			content := cleanHTML(item.Content)

			// Language detection
			if language, exists := u.detector.DetectLanguageOf(description + " " + content); exists {
				if language != lingua.English {
					stats.LanguageSkipped++
					log.Printf("Skipping non-English article from feed: %s", feedURL)
					continue
				}
			}

			// Extract categories
			category, subCategory := u.categorizer.Categorize(item, description, content)

			article := &models.Article{
				Link: item.Link,
			}

			if err := u.articleRepo.Create(ctx, article); err != nil {
				errMsg := fmt.Sprintf("Failed to save article %s: %v", item.Link, err)
				stats.ErrorDetails = append(stats.ErrorDetails, errMsg)
				stats.ErrorCount++
				log.Printf("Error saving article: %v", err)
				continue
			}

			message := command.NewArticleCommand{
				Title:       item.Title,
				Link:        item.Link,
				RSSLink:     feedURL,
				Description: description,
				Content:     content,
				Author:      getAuthor(item),
				PublishedAt: publishedAt,
				Category:    category,
				Subcategory: subCategory,
			}

			if err := u.sendToKafka(&message); err != nil {
				errMsg := fmt.Sprintf("Kafka delivery failed for %s: %v", item.Link, err)
				stats.ErrorDetails = append(stats.ErrorDetails, errMsg)
				stats.ErrorCount++
				log.Printf("Error sending to Kafka: %v", err)
				continue
			}

			stats.ProcessedCount++
		}

		stats.ProcessingTime = time.Since(feedStartTime)
		feedStats = append(feedStats, stats)
		totalCount += stats.ProcessedCount
	}

	totalProcessingTime := time.Since(startTime)

	// Send detailed Discord message
	if totalCount > 0 || len(globalErrors) > 0 {
		if err := u.sendDetailedDiscordMessage(feedStats, totalCount, totalProcessingTime, globalErrors); err != nil {
			log.Printf("Error sending Discord message: %v", err)
		}
	} else {
		log.Println("No new articles processed, skipping Discord notification.")
	}

	return totalCount, nil
}

func generateMessageChunks(stats []dto.FeedProcessingStats, totalProcessed int, totalTime time.Duration, errors []string) []string {
	var chunks []string
	const DISCORD_MAX_LENGTH = 2000

	var header strings.Builder
	header.WriteString(fmt.Sprintf("📰 **Feed Report** (finished at %s)\n", time.Now().Format("15:04:05")))
	header.WriteString(fmt.Sprintf("✅ **%d** processed in %d feeds, done in %vs \n", totalProcessed, len(stats), totalTime.Round(time.Second).Seconds()))

	currentChunk := header.String()

	var summaryBuilder strings.Builder
	for _, stat := range stats {
		if stat.ProcessedCount == 0 && stat.ErrorCount == 0 {
			continue
		}

		lineIcons := ""
		if stat.ErrorCount > 0 {
			lineIcons += fmt.Sprintf("   %d ⚠️", stat.ErrorCount)
		}
		if stat.SkippedCount > 0 {
			lineIcons += fmt.Sprintf("   %d ⏭️", stat.SkippedCount)
		}
		if stat.LanguageSkipped > 0 {
			lineIcons += fmt.Sprintf("   %d 🌍", stat.LanguageSkipped)
		}

		if lineIcons == "" {
			continue
		}

		line := fmt.Sprintf("• %s: %d ✅%s\n", stat.FeedURL, stat.ProcessedCount, lineIcons)

		if len(currentChunk)+len(line) > DISCORD_MAX_LENGTH {
			chunks = append(chunks, currentChunk)
			currentChunk = "..." + line
		} else {
			currentChunk += line
		}
	}

	if len(currentChunk) > len(header.String()) {
		summaryBuilder.WriteString(currentChunk)
		currentChunk = summaryBuilder.String()
	}

	if len(errors) > 0 {
		errorHeader := "\n⚠️ **Errors:**\n"

		if len(currentChunk)+len(errorHeader)+50 > DISCORD_MAX_LENGTH {
			chunks = append(chunks, currentChunk)
			currentChunk = ""
		}

		currentChunk += errorHeader

		for _, err := range errors {
			errorLine := fmt.Sprintf("• %s\n", err)

			if len(currentChunk)+len(errorLine) > DISCORD_MAX_LENGTH {
				chunks = append(chunks, currentChunk)
				currentChunk = "..." + errorLine
			} else {
				currentChunk += errorLine
			}
		}
	}

	if currentChunk != "" {
		chunks = append(chunks, currentChunk)
	}

	return chunks
}

func (u *RSSService) sendDetailedDiscordMessage(stats []dto.FeedProcessingStats, totalProcessed int, totalTime time.Duration, errors []string) error {
	if u.cfg.WebhookURL == "" {
		return nil
	}

	messageChunks := generateMessageChunks(stats, totalProcessed, totalTime, errors)

	username := "Albert - The News Bro"
	var combinedError error

	for _, chunk := range messageChunks {
		content := chunk

		err := discordwebhook.SendMessage(u.cfg.WebhookURL, discordwebhook.Message{
			Username: &username,
			Content:  &content,
		})

		if err != nil {
			if combinedError == nil {
				combinedError = err
			} else {
				combinedError = fmt.Errorf("%w; %w", combinedError, err)
			}
		}
	}

	return combinedError
}

func (u *RSSService) sendToKafka(message *command.NewArticleCommand) error {
	payload, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Key:   sarama.StringEncoder(message.Link),
		Topic: u.cfg.KafkaArticleCommandTopic,
		Value: sarama.ByteEncoder(payload),
	}

	_, _, err = u.producer.SendMessage(msg)
	return err
}

func (u *RSSService) HandleRSSAggregate(agg *aggregate.RSSAggregate) error {
	if agg.Active == false {
		return u.rssRepo.DeleteByLink(context.Background(), agg.Link)
	}
	exists, err := u.rssRepo.Exists(context.Background(), agg.Link)
	if err != nil {
		return err
	} else if exists {
		return nil
	}
	err = u.rssRepo.Create(context.Background(), &models.RSS{
		Link: agg.Link,
	})
	return err
}
