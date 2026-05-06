package broker

import (
	"encoding/json"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Event struct {
	EventID       string  `json:"event_id"`
	OrderID       int32   `json:"order_id"`
	Amount        float64 `json:"amount"`
	CustomerEmail string  `json:"customer_email"`
	Status        string  `json:"status"`
}

type Publisher struct {
	ch *amqp.Channel
}

func NewPublisher(ch *amqp.Channel) *Publisher {
	return &Publisher{ch: ch}
}

func (p *Publisher) Publish(event Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	err = p.ch.Publish(
		"",
		"payment.completed",
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		return err
	}

	log.Println("Published event:", string(body))
	return nil
}
