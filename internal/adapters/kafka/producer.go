package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/FollG/kafka-with-go/internal/domain/models"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/plain"
)

type Producer struct {
	writer     *kafka.Writer
	topic      string
	producerID string
}

func NewProducer(brokers []string, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        topic,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll, // exactly-once
		MaxAttempts:  3,
		BatchSize:    100,
		BatchTimeout: 10 * time.Millisecond,
		Async:        false,
	}

	return &Producer{
		writer:     writer,
		topic:      topic,
		producerID: fmt.Sprintf("producer-%d", time.Now().UnixNano()),
	}
}

func NewProducerWithAuth(brokers []string, topic, username, password string) *Producer {
	dialer := &kafka.Dialer{
		Timeout:   10 * time.Second,
		DualStack: true,
		SASLMechanism: plain.Mechanism{
			Username: username,
			Password: password,
		},
	}

	writer := kafka.WriterConfig{
		Brokers:  brokers,
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	}

	return &Producer{
		writer:     kafka.NewWriter(writer),
		topic:      topic,
		producerID: fmt.Sprintf("producer-%d", time.Now().UnixNano()),
	}
}

func (p *Producer) SendProductEvent(ctx context.Context, event *models.ProductEvent) error {
	event.ProducerID = p.producerID
	event.Sequence = time.Now().UnixNano()

	eventData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	msg := kafka.Message{
		Key:   []byte(fmt.Sprintf("product-%d", event.ProductID)),
		Value: eventData,
		Headers: []kafka.Header{
			{
				Key:   "event_type",
				Value: []byte(string(event.EventType)),
			},
			{
				Key:   "producer_id",
				Value: []byte(p.producerID),
			},
		},
		Time: time.Now(),
	}

	err = p.writer.WriteMessages(ctx, msg)
	if err != nil {
		return fmt.Errorf("failed to write message to kafka: %w", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
