package consumer

import (
	"context"
	"encoding/json"
	"log"

	"repo_feed/src/api/dto"
	"repo_feed/src/domain/services"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/command"
)

type InferenceConsumer struct {
	consumerGroup sarama.ConsumerGroup
	feedService   *services.FeedService
	topic         string
}

func NewInferenceConsumer(brokers []string, topic string, groupID string, feedService *services.FeedService) (*InferenceConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &InferenceConsumer{
		consumerGroup: consumerGroup,
		feedService:   feedService,
		topic:         topic,
	}, nil
}

func (ac *InferenceConsumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer for new recommendation...")
	handler := &consumerInferenceGroupHandler{feedService: ac.feedService}

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

func (ac *InferenceConsumer) Close() error {
	return ac.consumerGroup.Close()
}

type consumerInferenceGroupHandler struct {
	feedService *services.FeedService
}

func (h *consumerInferenceGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerInferenceGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerInferenceGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message received: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)
		h.processMessage(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerInferenceGroupHandler) processMessage(data []byte) {
	var cmd command.InferenceCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing new inference command for user ID: %d", cmd.UserID)

	article := dto.Article{
		ID:          cmd.Article.ID,
		Title:       cmd.Article.Title,
		Abstract:    cmd.Article.Abstract,
		Category:    cmd.Article.Category,
		Subcategory: cmd.Article.Subcategory,
		Link:        cmd.Article.Link,
		RSSLink:     cmd.Article.RSSLink,
		PublishedAt: cmd.Article.PublishedAt,
	}

	req := &dto.UpdateFeedRequest{
		UserID:  cmd.UserID,
		Article: &article,
		Model:   cmd.Model,
		Score:   cmd.Score,
	}

	err := h.feedService.UpdateFeedZSET(req)
	if err != nil {
		log.Printf("Error updating feed: %v", err)
		return
	}
}
