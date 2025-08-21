package main

import (
	"context"
	"encoding/json"
	"log"
	"ride-sharing/shared/contracts"
	"ride-sharing/shared/messaging"

	"github.com/rabbitmq/amqp091-go"
)

type tripEventConsumer struct {
	rabbitmq *messaging.RabbitMQ
	service  *Service
}

func NewTripEventConsumer(rabbitmq *messaging.RabbitMQ, service *Service) *tripEventConsumer {
	return &tripEventConsumer{
		rabbitmq: rabbitmq,
		service:  service,
	}
}

func (c *tripEventConsumer) Listen() error {
	return c.rabbitmq.ConsumeMessage(messaging.FindAvailableDriversQueue, func(ctx context.Context, msg amqp091.Delivery) error {
		var tripEvent contracts.AmqpMessage
		if err := json.Unmarshal(msg.Body, &tripEvent); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			return err
		}

		var payload messaging.TripEventData
		if err := json.Unmarshal(tripEvent.Data, &payload); err != nil {
			return err
		}

		switch msg.RoutingKey {
		case contracts.TripEventCreated, contracts.TripEventDriverNotInterested:
			return c.handleFindAndNotifyDrivers(ctx, payload)
		}

		log.Printf("unknown trip event: %+v", payload)

		return nil
	})
}

func (c *tripEventConsumer) handleFindAndNotifyDrivers(ctx context.Context, payload messaging.TripEventData) error {
	suitableIDs := c.service.FindAvailableDrivers(payload.Trip.SelectedFare.PackageSlug)

	log.Printf("Found suitable drivers %v", len(suitableIDs))

	if len(suitableIDs) == 0 {
		if err := c.rabbitmq.PublishMessage(ctx, contracts.TripEventNoDriversFound, contracts.AmqpMessage{
			OwnerID: payload.Trip.UserID,
		}); err != nil {
			log.Printf("Failed to publish message to exchange: %v", err)
		}

		return nil
	}

	suitableDriverID := suitableIDs[0]
	marshalledEvent, err := json.Marshal(suitableDriverID)
	if err != nil {
		return nil
	}

	if err := c.rabbitmq.PublishMessage(ctx, contracts.DriverCmdTripRequest, contracts.AmqpMessage{
		OwnerID: payload.Trip.UserID,
		Data:    marshalledEvent,
	}); err != nil {
		log.Printf("Failed to publish message to exchange: %v", err)
	}

	return nil
}
