package usecase

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type PaymentClient struct {
	client *http.Client
}

func NewPaymentClient() *PaymentClient {
	return &PaymentClient{
		client: &http.Client{
			Timeout: 2 * time.Second,
		},
	}
}

func (p *PaymentClient) Pay(orderID string, amount int64) (string, error) {
	payload := map[string]interface{}{
		"order_id": orderID,
		"amount":   amount,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	resp, err := p.client.Post(
		"http://localhost:8081/payments",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return "", errors.New("payment service unavailable")
	}
	defer resp.Body.Close()

	var result struct {
		Status string `json:"Status"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Status, nil
}
