package usecase

import (
	"errors"
	"order-service/internal/domain"
	"order-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo          repository.OrderRepository
	paymentClient *PaymentClient
}

func NewOrderUseCase(
	repo repository.OrderRepository,
	paymentClient *PaymentClient,
) *OrderUseCase {
	return &OrderUseCase{
		repo:          repo,
		paymentClient: paymentClient,
	}
}

func (u *OrderUseCase) Create(customerID, itemName string, amount int64) (*domain.Order, error) {
	if amount <= 0 {
		return nil, errors.New("amount must be > 0")
	}

	order := &domain.Order{
		ID:         uuid.New().String(),
		CustomerID: customerID,
		ItemName:   itemName,
		Amount:     amount,
		Status:     "Pending",
		CreatedAt:  time.Now(),
	}

	// 1. Save as pending
	if err := u.repo.Save(order); err != nil {
		return nil, err
	}

	// 2. Call payment service
	status, err := u.paymentClient.Pay(order.ID, order.Amount)
	if err != nil {
		// payment unavailable → failed
		_ = u.repo.UpdateStatus(order.ID, "Failed")
		return nil, errors.New("payment service unavailable")
	}

	// 3. Update order status
	if status == "Authorized" {
		order.Status = "Paid"
	} else {
		order.Status = "Failed"
	}

	if err := u.repo.UpdateStatus(order.ID, order.Status); err != nil {
		return nil, err
	}

	return order, nil
}

func (u *OrderUseCase) GetByID(id string) (*domain.Order, error) {
	return u.repo.GetByID(id)
}

func (u *OrderUseCase) Cancel(id string) error {
	order, err := u.repo.GetByID(id)
	if err != nil {
		return err
	}

	if order.Status != "Pending" {
		return errors.New("only pending orders can be cancelled")
	}

	return u.repo.UpdateStatus(id, "Cancelled")
}
