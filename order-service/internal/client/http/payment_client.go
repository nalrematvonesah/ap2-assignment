package http

import (
	"context"
	"time"

	paymentpb "github.com/nalrematvonesah/ap2-generated/paymentpb"
	"google.golang.org/grpc"
)

type PaymentClient struct {
	client paymentpb.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{
		client: paymentpb.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentClient) ProcessPayment(orderID int64, amount float64, userID string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.client.ProcessPayment(ctx, &paymentpb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
		UserId:  userID,
	})
	if err != nil {
		return "", err
	}

	return res.TransactionId, nil
}
