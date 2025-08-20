package main

import (
	"context"
	"log"
	"ride-sharing/shared/messaging"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type tripEventConsumer struct {
	rabbitmq *messaging.RabbitMQ
}

func NewTripEventConsumer(rabbitmq *messaging.RabbitMQ) *tripEventConsumer {
	return &tripEventConsumer{
		rabbitmq: rabbitmq,
	}
}

func (c *tripEventConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessage("hello", func(ctx context.Context, msg amqp091.Delivery) error {
		log.Printf("driver received message: %v", msg)
		time.Sleep(time.Second * 15)
		return nil
	})
}
