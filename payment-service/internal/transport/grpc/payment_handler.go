package grpc

import (
	"context"
	"fmt"

	paymentpb "github.com/nalrematvonesah/ap2-generated/paymentpb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"payment-service/internal/usecase"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) ProcessPayment(
	ctx context.Context,
	req *paymentpb.PaymentRequest,
) (*paymentpb.PaymentResponse, error) {

	orderID := fmt.Sprintf("%d", req.OrderId)
	amount := int64(req.Amount)

	payment, err := h.uc.Create(orderID, amount)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	success := payment.Status != "Declined"

	return &paymentpb.PaymentResponse{
		Success:       success,
		TransactionId: payment.TransactionID,
		Message:       payment.Status,
		ProcessedAt:   timestamppb.Now(),
	}, nil
}

func (h *PaymentHandler) ListPayments(
	ctx context.Context,
	req *paymentpb.ListPaymentsRequest,
) (*paymentpb.ListPaymentsResponse, error) {

	payments, err := h.uc.ListByStatus(req.Status)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var res []*paymentpb.PaymentResponse

	for _, p := range payments {
		res = append(res, &paymentpb.PaymentResponse{
			Success:       p.Status == "Authorized",
			TransactionId: p.TransactionID,
			Message:       p.Status,
			ProcessedAt:   timestamppb.Now(),
		})
	}

	return &paymentpb.ListPaymentsResponse{
		Payments: res,
	}, nil
}
