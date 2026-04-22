package usecase

import (
	"errors"
	grpcclient "order-service/internal/client/grpc"
	"order-service/internal/domain"
	"order-service/internal/repository"
	"time"

	"github.com/google/uuid"
)

type OrderUseCase struct {
	repo          repository.OrderRepository
	paymentClient *grpcclient.PaymentClient
}

func NewOrderUseCase(
	repo repository.OrderRepository,
	paymentClient *grpcclient.PaymentClient,
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

	// 1. Save
	if err := u.repo.Save(order); err != nil {
		return nil, err
	}

	// 2. gRPC call
	orderIDInt := time.Now().UnixNano()

	_, err := u.paymentClient.ProcessPayment(
		orderIDInt,
		float64(order.Amount),
		order.CustomerID,
	)

	if err != nil {
		_ = u.repo.UpdateStatus(order.ID, "Failed")
		return nil, err
	}

	// 3. status logic
	if order.Amount > 100000 {
		order.Status = "Failed"
	} else {
		order.Status = "Paid"
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
