package client

import (
	"context"
	"time"

	"google.golang.org/grpc"

	pb "github.com/nalrematvonesah/ap2-proto-contracts/gen/go/proto"
)

type PaymentClient struct {
	client pb.PaymentServiceClient
}

func NewPaymentClient(conn *grpc.ClientConn) *PaymentClient {
	return &PaymentClient{
		client: pb.NewPaymentServiceClient(conn),
	}
}

func (c *PaymentClient) ProcessPayment(orderID int32, amount float64, email string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := c.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: orderID,
		Amount:  amount,
		Email:   email,
	})

	if err != nil {
		return "", err
	}

	return res.Status, nil
}
