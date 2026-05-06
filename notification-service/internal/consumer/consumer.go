package consumer

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"

	"notification-service/internal/service"
)

type Event struct {
	EventID       string  `json:"event_id"`
	OrderID       int32   `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

type Consumer struct {
	service *service.NotificationService
	seen    map[string]bool
}

func NewConsumer(s *service.NotificationService) *Consumer {
	return &Consumer{
		service: s,
		seen:    make(map[string]bool),
	}
}

func (c *Consumer) Handle(msg amqp.Delivery) {
	var event Event

	err := json.Unmarshal(msg.Body, &event)
	if err != nil {
		log.Println("invalid message:", err)
		msg.Nack(false, false)
		return
	}

	if c.seen[event.EventID] {
		log.Println("duplicate event:", event.EventID)
		msg.Ack(false)
		return
	}

	c.seen[event.EventID] = true

	log.Println("Received event:", event)

	err = c.service.SendNotification(event)
	if err != nil {
		log.Println("failed to process:", err)
		msg.Nack(false, true) // retry
		return
	}

	msg.Ack(false)
}
