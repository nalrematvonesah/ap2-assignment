package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"notification-service/internal/consumer"
	"notification-service/internal/service"
)

func main() {
	var conn *amqp.Connection
	var err error

	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err == nil {
			break
		}

		log.Println("Retrying RabbitMQ connection...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(
		"payment.completed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("failed to consume: %v", err)
	}

	svc := service.NewNotificationService()
	c := consumer.NewConsumer(svc)

	log.Println("Notification Service started...")

	go func() {
		for msg := range msgs {
			c.Handle(msg)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down Notification Service...")

	log.Println("Closing RabbitMQ channel...")
	ch.Close()

	log.Println("Closing RabbitMQ connection...")
	conn.Close()

	log.Println("Notification Service stopped")
}
