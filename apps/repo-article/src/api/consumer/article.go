package consumer

import (
	"context"
	"encoding/json"

	"repo_article/src/api/dto"
	"repo_article/src/domain/services"

	"github.com/IBM/sarama"
	"github.com/TheRayquaza/newsbro/apps/libs/kafka/command"
	"github.com/TheRayquaza/newsbro/apps/libs/utils"
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
		utils.SugarLog.Errorf("failed to create consumer group: %s", err)
		return nil, err
	}

	return &ArticleConsumer{
		consumerGroup:  consumerGroup,
		articleService: articleService,
		topic:          topic,
	}, nil
}

func (ac *ArticleConsumer) Start(ctx context.Context) {
	utils.SugarLog.Infof("Starting Kafka consumer for new articles...")
	handler := &consumerGroupHandler{articleService: ac.articleService}

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
		utils.SugarLog.Infof("Message received: topic=%s, partition=%d, offset=%d", msg.Topic, msg.Partition, msg.Offset)
		h.processMessage(msg.Value)
		session.MarkMessage(msg, "")
	}
	return nil
}

func (h *consumerGroupHandler) processMessage(data []byte) {
	var cmd command.NewArticleCommand
	if err := json.Unmarshal(data, &cmd); err != nil {
		utils.SugarLog.Errorf("Error unmarshaling message: %v", err)
		return
	}

	utils.SugarLog.Infof("Processing new article: %s", cmd.Title)

	req := &dto.ArticleCreateRequest{
		Title:       cmd.Title,
		Link:        cmd.Link,
		RSSLink:     cmd.RSSLink,
		Abstract:    cmd.Description,
		Category:    cmd.Category,
		Subcategory: cmd.Subcategory,
		PublishedAt: cmd.PublishedAt,
	}

	article, err := h.articleService.CreateArticle(req, 0)
	if err != nil {
		utils.SugarLog.Errorf("Error creating article: %v", err)
		return
	}

	utils.SugarLog.Infof("Successfully created article with ID: %d", article.ID)
}
