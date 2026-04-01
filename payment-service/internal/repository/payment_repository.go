package repository

import (
	"database/sql"
	"payment-service/internal/domain"
)

type PaymentRepository interface {
	Save(payment *domain.Payment) error
	GetByOrderID(orderID string) (*domain.Payment, error)
}

type paymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

func (r *paymentRepository) Save(payment *domain.Payment) error {
	_, err := r.db.Exec(`
		INSERT INTO payments(id, order_id, transaction_id, amount, status)
		VALUES($1,$2,$3,$4,$5)
	`, payment.ID, payment.OrderID, payment.TransactionID, payment.Amount, payment.Status)
	return err
}

func (r *paymentRepository) GetByOrderID(orderID string) (*domain.Payment, error) {
	row := r.db.QueryRow(`
		SELECT id, order_id, transaction_id, amount, status
		FROM payments WHERE order_id=$1
	`, orderID)

	var p domain.Payment
	err := row.Scan(&p.ID, &p.OrderID, &p.TransactionID, &p.Amount, &p.Status)
	return &p, err
}
