package kafka

import (
	"fmt"

	"github.com/segmentio/kafka-go"
)

func CreateTopic(brokers []string, topic string, partitions int) error {
	conn, err := kafka.Dial("tcp", brokers[0])
	if err != nil {
		return fmt.Errorf("failed to dial kafka: %w", err)
	}
	defer func(conn *kafka.Conn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: 1,
		},
	}

	err = conn.CreateTopics(topicConfigs...)
	if err != nil {
		// Топик может уже существовать - это нормально
		fmt.Printf("Note: topic creation may have failed (it might already exist): %v\n", err)
	}

	return nil
}
