package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	amqp "github.com/rabbitmq/amqp091-go"

	pb "github.com/nalrematvonesah/ap2-proto-contracts/gen/go/proto"

	"payment-service/internal/broker"
	"payment-service/internal/database"
	"payment-service/internal/handler"
	"payment-service/internal/service"
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
		log.Fatalf("failed to connect to RabbitMQ: %v", err)
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

	db := database.NewPostgres()

	publisher := broker.NewPublisher(ch)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	svc := service.NewPaymentService(db)

	server := handler.NewServer(svc, publisher)

	grpcServer := grpc.NewServer()

	pb.RegisterPaymentServiceServer(grpcServer, server)
	reflection.Register(grpcServer)

	go func() {
		log.Println("Payment Service running on :50051")

		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(
		stop,
		os.Interrupt,
		syscall.SIGTERM,
	)

	<-stop

	log.Println("Shutting down Payment Service...")

	grpcServer.GracefulStop()

	log.Println("Closing DB...")
	db.Close()

	log.Println("Closing RabbitMQ channel...")
	ch.Close()

	log.Println("Closing RabbitMQ connection...")
	conn.Close()

	log.Println("Payment Service stopped")
}
