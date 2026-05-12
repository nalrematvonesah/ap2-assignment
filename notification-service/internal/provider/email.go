package provider

import (
	"errors"
	"log"
	"math/rand"
	"time"
)

type EmailSender interface {
	Send(to string, body string) error
}

type MockEmailSender struct{}

func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

func (m *MockEmailSender) Send(
	to string,
	body string,
) error {

	time.Sleep(2 * time.Second)

	if rand.Intn(100) < 40 {
		return errors.New("simulated provider failure")
	}

	log.Println("EMAIL SENT TO:", to)

	return nil
}
