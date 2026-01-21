package shared

type KafkaConfig struct {
	Topic         string
	ConsumerGroup string
	Host          string
}

func NewKafkaConfig() *KafkaConfig {
	return &KafkaConfig{
		Topic:         "example-topic",
		ConsumerGroup: "example-group",
		Host:          "localhost:9092",
	}
}
