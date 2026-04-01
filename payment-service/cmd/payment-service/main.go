package main

import (
	"database/sql"
	"log"
	"payment-service/internal/app"
	"payment-service/internal/repository"
	paymentHttp "payment-service/internal/transport/http"
	"payment-service/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	handler := paymentHttp.NewPaymentHandler(uc)

	r := gin.Default()

	r.Use(cors.Default())

	app.SetupRoutes(r, handler)

	log.Println("Payment Service started on :8081")
	if err := r.Run(":8081"); err != nil {
		log.Fatal(err)
	}
}
