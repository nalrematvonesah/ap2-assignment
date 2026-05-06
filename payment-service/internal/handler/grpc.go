package handler

import (
	"context"

	"github.com/google/uuid"

	pb "github.com/nalrematvonesah/ap2-proto-contracts/gen/go/proto"

	"payment-service/internal/broker"
	"payment-service/internal/service"
)

type Server struct {
	pb.UnimplementedPaymentServiceServer
	service   *service.PaymentService
	publisher *broker.Publisher
}

func NewServer(s *service.PaymentService, p *broker.Publisher) *Server {
	return &Server{
		service:   s,
		publisher: p,
	}
}

func (s *Server) ProcessPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	status := s.service.Process(req.OrderId, req.Amount, req.Email)

	event := broker.Event{
		EventID:       uuid.New().String(),
		OrderID:       req.OrderId,
		Amount:        req.Amount,
		CustomerEmail: req.Email,
		Status:        status,
	}

	err := s.publisher.Publish(event)
	if err != nil {
		return nil, err
	}

	return &pb.PaymentResponse{
		Status: status,
	}, nil
}
