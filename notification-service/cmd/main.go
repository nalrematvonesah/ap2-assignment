package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"

	"notification-service/internal/consumer"
	"notification-service/internal/provider"
	"notification-service/internal/service"
)

func main() {
	var conn *amqp.Connection
	var err error

	for {
		conn, err = amqp.Dial(
			"amqp://guest:guest@rabbitmq:5672/",
		)

		if err == nil {
			log.Println(
				"Connected to RabbitMQ",
			)
			break
		}

		log.Println(
			"Retrying RabbitMQ connection...",
		)

		time.Sleep(3 * time.Second)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf(
			"failed to open channel: %v",
			err,
		)
	}

	defer ch.Close()

	args := amqp.Table{
		"x-dead-letter-exchange":    "",
		"x-dead-letter-routing-key": "payment.completed.dlq",
	}

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		args,
	)

	if err != nil {
		log.Fatalf(
			"failed to declare queue: %v",
			err,
		)
	}

	_, err = ch.QueueDeclare(
		"payment.completed.dlq",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Fatalf(
			"failed to declare DLQ: %v",
			err,
		)
	}

	err = ch.Qos(
		1,
		0,
		false,
	)

	if err != nil {
		log.Fatalf(
			"failed to set qos: %v",
			err,
		)
	}

	var redisClient *redis.Client

	for i := 0; i < 10; i++ {

		redisClient = redis.NewClient(
			&redis.Options{
				Addr: "redis:6379",
			},
		)

		_, err = redisClient.Ping(
			context.Background(),
		).Result()

		if err == nil {
			log.Println(
				"Connected to Redis",
			)
			break
		}

		log.Println(
			"Retrying Redis connection...",
		)

		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf(
			"failed to connect redis: %v",
			err,
		)
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
		log.Fatalf(
			"failed to consume: %v",
			err,
		)
	}

	sender := provider.NewMockEmailSender()

	svc := service.NewNotificationService(
		sender,
	)

	c := consumer.NewConsumer(
		svc,
		redisClient,
	)

	log.Println(
		"Notification Worker started...",
	)

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

	log.Println(
		"Shutting down Notification Worker...",
	)

	log.Println("Closing Redis...")
	redisClient.Close()

	log.Println(
		"Closing RabbitMQ channel...",
	)
	ch.Close()

	log.Println(
		"Closing RabbitMQ connection...",
	)
	conn.Close()

	log.Println(
		"Notification Worker stopped",
	)
}
