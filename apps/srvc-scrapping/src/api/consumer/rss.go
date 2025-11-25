package consumer

import (
	"context"
	"encoding/json"
	"time"

	"srvc_scrapping/src/domain/service"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/aggregate"
	"github.com/TheRayquaza/newsbro/apps/libs/utils"
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
	utils.SugarLog.Infof("Starting Kafka consumer for new rsss...")
	handler := &consumerGroupHandler{rssService: ac.rssService}

	for {
		select {
		case <-ctx.Done():
			utils.SugarLog.Infof("Shutting down Kafka consumer...")
			if err := ac.consumerGroup.Close(); err != nil {
				utils.SugarLog.Errorf("Error closing consumer group: %v", err)
			}
			return
		default:
			if err := ac.consumerGroup.Consume(ctx, []string{ac.topic}, handler); err != nil {
				utils.SugarLog.Errorf("Error from consumer: %v", err)
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
		utils.SugarLog.Infof("Message received: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)
		h.processMessage(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerGroupHandler) processMessage(data []byte) {
	var cmd aggregate.RSSAggregate
	if err := json.Unmarshal(data, &cmd); err != nil {
		utils.SugarLog.Errorf("Error unmarshaling message: %v", err)
		return
	}

	if err := h.rssService.HandleRSSAggregate(&cmd); err != nil {
		utils.SugarLog.Errorf("Error processing new rss: %v", err)
	}
}
