package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"srvc_scrapping/src/domain/service"
)

type RSSConsumer struct {
	consumerGroup sarama.ConsumerGroup
	rssService    *service.RSSService
	topic         string
}

func NewRSSConsumer(brokers []string, topic string, groupID string, rssService *service.RSSService) (*RSSConsumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRoundRobin()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}

	return &RSSConsumer{
		consumerGroup: consumerGroup,
		rssService:    rssService,
		topic:         topic,
	}, nil
}

func (ac *RSSConsumer) Start(ctx context.Context) {
	log.Println("Starting Kafka consumer for new rsss...")
	handler := &consumerGroupHandler{rssService: ac.rssService}

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

func (ac *RSSConsumer) Close() error {
	return ac.consumerGroup.Close()
}

type consumerGroupHandler struct {
	rssService *service.RSSService
}

func (h *consumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *consumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (ac *RSSConsumer) ConsumeUntilEmpty(ctx context.Context, timeout time.Duration) error {
	handler := &consumerGroupHandler{rssService: ac.rssService}

	timeoutCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	errChan := make(chan error, 1)

	go func() {
		errChan <- ac.consumerGroup.Consume(timeoutCtx, []string{ac.topic}, handler)
	}()

	select {
	case err := <-errChan:
		return err
	case <-timeoutCtx.Done():
		return nil
	}
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
	var cmd aggregate.RSSAggregate
	if err := json.Unmarshal(data, &cmd); err != nil {
		log.Printf("Error unmarshaling message: %v", err)
		return
	}

	if err := h.rssService.HandleRSSAggregate(&cmd); err != nil {
		log.Printf("Error processing new rss: %v", err)
	}
}
