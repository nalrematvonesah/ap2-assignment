package service

import (
	"fmt"

	"notification-service/internal/provider"
)

type Event struct {
	EventID       string
	OrderID       int32
	Amount        float64
	CustomerEmail string
	Status        string
}

type NotificationService struct {
	sender provider.EmailSender
}

func NewNotificationService(
	sender provider.EmailSender,
) *NotificationService {

	return &NotificationService{
		sender: sender,
	}
}

func (s *NotificationService) SendNotification(
	e Event,
) error {

	body := fmt.Sprintf(
		"Payment completed for order %d amount %.2f",
		e.OrderID,
		e.Amount,
	)

	return s.sender.Send(
		e.CustomerEmail,
		body,
	)
}
