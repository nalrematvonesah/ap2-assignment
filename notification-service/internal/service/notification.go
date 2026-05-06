package service

import "log"

type Event struct {
	EventID       string
	OrderID       int32
	Amount        float64
	CustomerEmail string
	Status        string
}

type NotificationService struct{}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(e interface{}) error {
	log.Println("Sending notification:", e)
	return nil
}
