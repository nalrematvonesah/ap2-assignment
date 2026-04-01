package repository

import (
	"database/sql"
	"order-service/internal/domain"
)

type OrderRepository interface {
	Save(order *domain.Order) error
	UpdateStatus(id string, status string) error
	GetByID(id string) (*domain.Order, error)
}

type orderRepository struct {
	db *sql.DB
}

func NewOrderRepository(db *sql.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) Save(order *domain.Order) error {
	_, err := r.db.Exec(`
		INSERT INTO orders(id, customer_id, item_name, amount, status, created_at)
		VALUES($1,$2,$3,$4,$5,$6)
	`, order.ID, order.CustomerID, order.ItemName, order.Amount, order.Status, order.CreatedAt)
	return err
}

func (r *orderRepository) UpdateStatus(id string, status string) error {
	_, err := r.db.Exec(`
		UPDATE orders SET status=$1 WHERE id=$2
	`, status, id)
	return err
}

func (r *orderRepository) GetByID(id string) (*domain.Order, error) {
	row := r.db.QueryRow(`
		SELECT id, customer_id, item_name, amount, status, created_at
		FROM orders WHERE id=$1
	`, id)

	var order domain.Order
	err := row.Scan(
		&order.ID,
		&order.CustomerID,
		&order.ItemName,
		&order.Amount,
		&order.Status,
		&order.CreatedAt,
	)

	return &order, err
}
