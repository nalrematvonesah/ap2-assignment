package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type Payment struct {
	OrderID int32   `json:"order_id"`
	Amount  float64 `json:"amount"`
	Email   string  `json:"email"`
	Status  string  `json:"status"`
}

type PaymentService struct {
	db    *sql.DB
	redis *redis.Client
}

func NewPaymentService(
	db *sql.DB,
	redis *redis.Client,
) *PaymentService {

	return &PaymentService{
		db:    db,
		redis: redis,
	}
}

func (s *PaymentService) Process(
	orderID int32,
	amount float64,
	email string,
) string {

	log.Printf(
		"Processing payment: order=%d amount=%.2f email=%s",
		orderID,
		amount,
		email,
	)

	_, err := s.db.Exec(
		`
		INSERT INTO payments(
			order_id,
			amount,
			email,
			status
		)
		VALUES($1,$2,$3,$4)
		`,
		orderID,
		amount,
		email,
		"completed",
	)

	if err != nil {
		log.Println("DB insert error:", err)
	}

	payment := Payment{
		OrderID: orderID,
		Amount:  amount,
		Email:   email,
		Status:  "completed",
	}

	data, _ := json.Marshal(payment)

	key := fmt.Sprintf(
		"payment:%d",
		orderID,
	)

	s.redis.Set(
		context.Background(),
		key,
		data,
		10*time.Minute,
	)

	return "completed"
}

func (s *PaymentService) GetPayment(
	orderID int32,
) (*Payment, error) {

	ctx := context.Background()

	key := fmt.Sprintf(
		"payment:%d",
		orderID,
	)

	cached, err := s.redis.Get(
		ctx,
		key,
	).Result()

	if err == nil {
		log.Println("CACHE HIT")

		var payment Payment

		json.Unmarshal(
			[]byte(cached),
			&payment,
		)

		return &payment, nil
	}

	log.Println("CACHE MISS")

	row := s.db.QueryRow(
		`
		SELECT
			order_id,
			amount,
			email,
			status
		FROM payments
		WHERE order_id = $1
		`,
		orderID,
	)

	var payment Payment

	err = row.Scan(
		&payment.OrderID,
		&payment.Amount,
		&payment.Email,
		&payment.Status,
	)

	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(payment)

	s.redis.Set(
		ctx,
		key,
		data,
		10*time.Minute,
	)

	return &payment, nil
}
