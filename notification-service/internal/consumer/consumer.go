package consumer

import (
	"context"
	"encoding/json"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

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
	redis   *redis.Client
}

func NewConsumer(
	s *service.NotificationService,
	r *redis.Client,
) *Consumer {

	return &Consumer{
		service: s,
		redis:   r,
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

	ctx := context.Background()

	key := "event:" + event.EventID

	exists, _ := c.redis.Exists(ctx, key).Result()

	if exists == 1 {
		log.Println("duplicate event:", event.EventID)
		msg.Ack(false)
		return
	}

	log.Println("Received event:", event)

	var processErr error

	for attempt := 1; attempt <= 3; attempt++ {
		processErr = c.service.SendNotification(
			service.Event(event),
		)

		if processErr == nil {
			break
		}

		log.Println(
			"retry attempt:",
			attempt,
		)

		time.Sleep(
			time.Duration(attempt*2) * time.Second,
		)
	}

	if processErr != nil {
		log.Println(
			"failed after retries:",
			processErr,
		)

		msg.Nack(false, false)
		return
	}

	c.redis.Set(
		ctx,
		key,
		"processed",
		24*time.Hour,
	)

	msg.Ack(false)
}
