package main

import (
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	"gateway/internal/handler"
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

		log.Println("Retrying gRPC connection...")
		time.Sleep(3 * time.Second)
	}

	if err != nil {
		log.Fatalf("failed to connect to payment service: %v", err)
	}

	h := handler.NewPaymentHandler(conn)

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:3000",
		},
		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"DELETE",
			"OPTIONS",
		},
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
		},
	}))

	r.POST("/payments", h.ProcessPayment)

	log.Println("Gateway running on :8080")

	r.Run(":8080")
}
