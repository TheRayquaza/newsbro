package consumer

import (
	"context"
	"encoding/json"
	"log"

	"repo_feed/src/domain/services"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
)

type FeedbackConsumer struct {
	consumerGroup sarama.ConsumerGroup
	feedService   *services.FeedService
	topic         string
}

func NewFeedbackConsumer(brokers []string, topic string, groupID string, feedService *services.FeedService) (*FeedbackConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &FeedbackConsumer{
		consumerGroup: consumerGroup,
		feedService:   feedService,
		topic:         topic,
	}, nil
}

func (ac *FeedbackConsumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer for new feedbacks...")
	handler := &consumerFeedbackGroupHandler{feedService: ac.feedService}

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

func (ac *FeedbackConsumer) Close() error {
	return ac.consumerGroup.Close()
}

type consumerFeedbackGroupHandler struct {
	feedService *services.FeedService
}

func (h *consumerFeedbackGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerFeedbackGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerFeedbackGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.Printf("Message received: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)
		h.processMessage(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerFeedbackGroupHandler) processMessage(data []byte) {
	var cmd aggregate.FeedbackAggregate
	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	log.Printf("Processing new feedback for user %d on article %d", cmd.UserID, cmd.NewsID)

	err := h.feedService.AddNewFeedbackAllModels(cmd.UserID, cmd.NewsID)
	if err != nil {
		log.Printf("Error updating feed: %v", err)
		return
	}
}
