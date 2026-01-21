package producer

import (
	"fmt"
	"kafka/internal/shared"

	kafka "github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

func NewKafkaProducer(kafkaConfig *shared.KafkaConfig) (*KafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaConfig.Host,
	})
	if err != nil {
		return nil, err
	}

	kp := &KafkaProducer{
		producer: p,
		topic:    kafkaConfig.Topic,
	}

	// Listen for delivery reports
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	return kp, nil
}

// Produce a message
func (kp *KafkaProducer) Produce(message string) error {
	return kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &kp.topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(message),
	}, nil)
}

// Close the producer gracefully
func (kp *KafkaProducer) Close() {
	kp.producer.Flush(15000) // wait max 15s for messages to deliver
	kp.producer.Close()
}
