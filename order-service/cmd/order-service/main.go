package main

import (
	"database/sql"
	"log"
	"order-service/internal/app"
	"order-service/internal/repository"
	orderHttp "order-service/internal/transport/http"
	"order-service/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	dsn := "postgres://postgres:postgres@127.0.0.1:55432/order_db?sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("sql.Open error:", err)
	}

	// 🔥 THIS IS IMPORTANT
	if err := db.Ping(); err != nil {
		log.Fatal("db.Ping error:", err)
	}

	log.Println("connected to order_db")

	repo := repository.NewOrderRepository(db)
	paymentClient := usecase.NewPaymentClient()
	uc := usecase.NewOrderUseCase(repo, paymentClient)
	handler := orderHttp.NewOrderHandler(uc)

	r := gin.Default()
	r.Use(cors.Default())

	app.SetupRoutes(r, handler)

	log.Println("Order Service started on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
