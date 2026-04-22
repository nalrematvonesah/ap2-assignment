package main

import (
	"database/sql"
	"log"
	"net"
	"order-service/internal/app"
	grpcclient "order-service/internal/client/grpc"
	"order-service/internal/repository"
	grpcTransport "order-service/internal/transport/grpc"
	http "order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	orderpb "github.com/nalrematvonesah/ap2-generated/orderpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	dsn := "postgres://postgres:postgres@127.0.0.1:55432/order_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// =====================
	// DB + repo
	// =====================
	repo := repository.NewOrderRepository(db)

	// =====================
	// gRPC client (Payment)
	// =====================
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	paymentClient := grpcclient.NewPaymentClient(conn)

	// =====================
	// Usecase
	// =====================
	uc := usecase.NewOrderUseCase(repo, paymentClient)

	// =====================
	// REST (Gin)
	// =====================
	handler := http.NewOrderHandler(uc)

	r := gin.Default()
	r.Use(cors.Default())

	app.SetupRoutes(r, handler)

	go func() {
		log.Println("REST Order Service started on :8080")
		if err := r.Run(":8080"); err != nil {
			log.Fatal(err)
		}
	}()

	// =====================
	// gRPC server (Streaming)
	// =====================
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(grpcTransport.LoggingUnaryInterceptor),
		grpc.StreamInterceptor(grpcTransport.LoggingStreamInterceptor),
	)

	grpcHandler := grpcTransport.NewOrderHandler(uc)

	orderpb.RegisterOrderServiceServer(grpcServer, grpcHandler)

	log.Println("gRPC Order Service running on :50052")
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
