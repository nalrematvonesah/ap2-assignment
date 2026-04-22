package main

import (
	"database/sql"
	"log"
	"net"
	"os"
	"payment-service/internal/app"
	"payment-service/internal/repository"
	grpcTransport "payment-service/internal/transport/grpc"
	paymentHttp "payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"google.golang.org/grpc/reflection"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	paymentpb "github.com/nalrematvonesah/ap2-generated/paymentpb"
	"google.golang.org/grpc"
)

func main() {
	dsn := "postgres://postgres:postgres@127.0.0.1:55433/payment_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("db.Ping error:", err)
	}

	repo := repository.NewPaymentRepository(db)
	uc := usecase.NewPaymentUseCase(repo)

	// =====================
	// REST (как было)
	// =====================
	httpHandler := paymentHttp.NewPaymentHandler(uc)

	r := gin.Default()
	r.Use(cors.Default())
	app.SetupRoutes(r, httpHandler)

	go func() {
		log.Println("REST Payment Service started on :8081")
		if err := r.Run(":8081"); err != nil {
			log.Fatal(err)
		}
	}()

	// =====================
	// gRPC
	// =====================
	port := os.Getenv("PAYMENT_GRPC_PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal(err)
	}

	grpcHandler := grpcTransport.NewPaymentHandler(uc)
	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, grpcHandler)
	reflection.Register(grpcServer)
	log.Println("gRPC Payment Service running on :" + port)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
