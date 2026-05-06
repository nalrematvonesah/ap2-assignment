package main

import (
	"log"
	"time"

	"google.golang.org/grpc"

	"order-service/internal/client"
)

func main() {
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 10; i++ {
		conn, err = grpc.Dial(
			"payment-service:50051",
			grpc.WithInsecure(),
			grpc.WithBlock(),
		)

		if err == nil {
			break
		}

		log.Println("Retrying connection to Payment Service...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := client.NewPaymentClient(conn)

	status, err := client.ProcessPayment(1, 150.0, "user@test.com")
	if err != nil {
		log.Fatalf("error calling payment service: %v", err)
	}

	log.Println("Payment status:", status)
}
