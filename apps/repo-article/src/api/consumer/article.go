package consumer

import (
	"context"
	"encoding/json"
	"log"

	"repo_article/src/api/dto"
	"repo_article/src/domain/services"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/command"
)

type ArticleConsumer struct {
	consumerGroup  sarama.ConsumerGroup
	articleService *services.ArticleService
	topic          string
}

func NewArticleConsumer(brokers []string, topic string, groupID string, articleService *services.ArticleService) (*ArticleConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &ArticleConsumer{
		consumerGroup:  consumerGroup,
		articleService: articleService,
		topic:          topic,
	}, nil
}

func (ac *ArticleConsumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer for new articles...")
	handler := &consumerGroupHandler{articleService: ac.articleService}

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down Kafka consumer...")
			if err := ac.consumerGroup.Close(); err != nil {
				log.Printf("Error closing consumer group: %v", err)
			}
			return
		default:
			if err := ac.consumerGroup.Consume(ctx, []string{ac.topic}, handler); err != nil {
				log.Printf("Error from consumer: %v", err)
			}
		}
	}
}

func (ac *ArticleConsumer) Close() error {
	return ac.consumerGroup.Close()
}

type consumerGroupHandler struct {
	articleService *services.ArticleService
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message received: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)
		h.processMessage(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerGroupHandler) processMessage(data []byte) {
	var cmd command.NewArticleCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing new article: %s", cmd.Title)

	req := &dto.ArticleCreateRequest{
		Title:       cmd.Title,
		Link:        cmd.Link,
		Abstract:    cmd.Description,
		Category:    cmd.Category,
		Subcategory: cmd.Subcategory,
	}

	article, err := h.articleService.CreateArticle(req, 0)
	if err != nil {
		log.Printf("Error creating article: %v", err)
		return
	}

	log.Printf("Successfully created article with ID: %d", article.ID)
}
