package service

import (
	"database/sql"
	"log"
)

type PaymentService struct {
	db *sql.DB
}

func NewPaymentService(db *sql.DB) *PaymentService {
	return &PaymentService{
		db: db,
	}
}

func (s *PaymentService) Process(orderID int32, amount float64, email string) string {
	log.Printf(
		"Processing payment: order=%d amount=%.2f email=%s",
		orderID,
		amount,
		email,
	)

	_, err := s.db.Exec(
		"INSERT INTO payments(order_id, amount, email, status) VALUES($1,$2,$3,$4)",
		orderID,
		amount,
		email,
		"completed",
	)

	if err != nil {
		log.Println("DB insert error:", err)
	}

	return "completed"
}
