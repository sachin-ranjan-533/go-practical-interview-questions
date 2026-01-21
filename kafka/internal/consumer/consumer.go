package consumer

import (
	"fmt"
	"kafka/internal/shared"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaConsumer struct {
	consumer *kafka.Consumer
	topic    string
	msgCH    chan<- string
}

func NewKafkaConsumer(kafkaConfig *shared.KafkaConfig, msgCH chan<- string) (*KafkaConsumer, error) {
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Host,
		"group.id":          kafkaConfig.ConsumerGroup,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer: c,
		topic:    kafkaConfig.Topic,
		msgCH:    msgCH,
	}, nil
}

// Consume messages continuously
func (kc *KafkaConsumer) Consume() error {
	// Subscribe to topic
	if err := kc.consumer.SubscribeTopics([]string{kc.topic}, nil); err != nil {
		return err
	}

	for {
		msg, err := kc.consumer.ReadMessage(time.Second)
		if err == nil {
			payload := fmt.Sprintf("Message on %s: %s", msg.TopicPartition, string(msg.Value))
			kc.msgCH <- payload
		} else if !err.(kafka.Error).IsTimeout() {
			fmt.Printf("Consumer error: %v\n", err)
		}
	}

	return nil
}

// Optional: Close consumer gracefully
func (kc *KafkaConsumer) Close() {
	kc.consumer.Close()
}
