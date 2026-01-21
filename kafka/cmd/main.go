package main

import (
	"fmt"
	"kafka/internal/consumer"
	"kafka/internal/producer"
	"kafka/internal/shared"
	"time"
)

type Server struct {
	producer *producer.KafkaProducer
	consumer *consumer.KafkaConsumer
}

func NewServer() (*Server, error) {
	msgCH := make(chan string, 10)

	// Create producer
	p, err := producer.NewKafkaProducer(shared.NewKafkaConfig())
	if err != nil {
		return nil, err
	}

	// Create consumer
	c, err := consumer.NewKafkaConsumer(shared.NewKafkaConfig(), msgCH)
	if err != nil {
		return nil, err
	}

	// Print consumed messages
	go func() {
		for msg := range msgCH {
			fmt.Println(msg)
		}
	}()

	return &Server{
		producer: p,
		consumer: c,
	}, nil
}

// Send messages every 2 seconds
func (s *Server) SendMessage() {
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	id := 0
	for range ticker.C {
		msg := fmt.Sprintf("message from kafka server = %d", id)
		id++

		if err := s.producer.Produce(msg); err != nil {
			fmt.Println("Producer not ready, retrying:", err)
			continue
		}

		fmt.Println("Sent:", msg)
	}
}

func (s *Server) ReceiveMessage() error {
	return s.consumer.Consume()
}

func main() {
	server, err := NewServer()
	if err != nil {
		panic(err)
	}

	// Run producer in a separate goroutine
	go server.SendMessage()

	// Run consumer in main goroutine
	if err := server.ReceiveMessage(); err != nil {
		panic(err)
	}
}
