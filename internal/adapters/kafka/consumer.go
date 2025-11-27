package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FollG/kafka-with-go/internal/domain/models"
	"github.com/FollG/kafka-with-go/internal/domain/repositories"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Consumer struct {
	reader       *kafka.Reader
	productRepo  repositories.ProductRepository
	cache        repositories.ProductCache
	batchSize    int
	batchTimeout time.Duration
}

func NewConsumer(brokers []string, topic, groupID string, batchSize int, batchTimeout time.Duration) *Consumer {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3, // 10KB
		MaxBytes:       10e6, // 10MB
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		MaxWait:        batchTimeout,
		QueueCapacity:  batchSize,
	})

	return &Consumer{
		reader:       reader,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}
}

func NewConsumerWithAuth(brokers []string, topic, groupID, username, password string, batchSize int, batchTimeout time.Duration) *Consumer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: username,
			Password: password,
		},
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        brokers,
		Topic:          topic,
		GroupID:        groupID,
		MinBytes:       10e3,
		MaxBytes:       10e6,
		CommitInterval: time.Second,
		StartOffset:    kafka.LastOffset,
		MaxWait:        batchTimeout,
		QueueCapacity:  batchSize,
		Dialer:         dialer,
	})

	return &Consumer{
		reader:       reader,
		batchSize:    batchSize,
		batchTimeout: batchTimeout,
	}
}

func (c *Consumer) SetProductRepo(repo repositories.ProductRepository) {
	c.productRepo = repo
}

func (c *Consumer) SetCache(cache repositories.ProductCache) {
	c.cache = cache
}

func (c *Consumer) Start(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg, err := c.reader.FetchMessage(ctx)
			if err != nil {
				return fmt.Errorf("failed to fetch message: %w", err)
			}

			if err := c.processMessage(ctx, msg); err != nil {
				fmt.Printf("Failed to process message: %v\n", err)
				continue
			}

			if err := c.reader.CommitMessages(ctx, msg); err != nil {
				return fmt.Errorf("failed to commit message: %w", err)
			}
		}
	}
}

func (c *Consumer) processMessage(ctx context.Context, msg kafka.Message) error {
	var event models.ProductEvent
	if err := json.Unmarshal(msg.Value, &event); err != nil {
		return fmt.Errorf("failed to unmarshal event: %w", err)
	}

	switch event.EventType {
	case models.ProductCreated:
		return c.handleProductCreated(ctx, &event)
	case models.ProductUpdated:
		return c.handleProductUpdated(ctx, &event)
	case models.ProductDeleted:
		return c.handleProductDeleted(ctx, &event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}

func (c *Consumer) handleProductCreated(ctx context.Context, event *models.ProductEvent) error {
	if event.ProductData == nil {
		return fmt.Errorf("product data is nil for create event")
	}

	if err := c.productRepo.Create(ctx, event.ProductData); err != nil {
		return fmt.Errorf("failed to create product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", event.ProductData.ID)
	if err := c.cache.Set(ctx, cacheKey, event.ProductData); err != nil {
		fmt.Printf("Failed to cache product %d: %v\n", event.ProductData.ID, err)
	}

	fmt.Printf("Successfully created product: %d\n", event.ProductData.ID)
	return nil
}

func (c *Consumer) handleProductUpdated(ctx context.Context, event *models.ProductEvent) error {
	if event.ProductData == nil {
		return fmt.Errorf("product data is nil for update event")
	}

	if err := c.productRepo.Update(ctx, event.ProductData); err != nil {
		return fmt.Errorf("failed to update product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", event.ProductData.ID)
	if err := c.cache.Set(ctx, cacheKey, event.ProductData); err != nil {
		fmt.Printf("Failed to cache product %d: %v\n", event.ProductData.ID, err)
	}

	fmt.Printf("Successfully updated product: %d\n", event.ProductData.ID)
	return nil
}

func (c *Consumer) handleProductDeleted(ctx context.Context, event *models.ProductEvent) error {
	if err := c.productRepo.Delete(ctx, event.ProductID); err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	cacheKey := fmt.Sprintf("product:%d", event.ProductID)
	if err := c.cache.Delete(ctx, cacheKey); err != nil {
		fmt.Printf("Failed to delete product %d from cache: %v\n", event.ProductID, err)
	}

	fmt.Printf("Successfully deleted product: %d\n", event.ProductID)
	return nil
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
