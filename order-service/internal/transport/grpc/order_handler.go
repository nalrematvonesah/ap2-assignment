package grpc

import (
	"fmt"
	"time"

	orderpb "github.com/nalrematvonesah/ap2-generated/orderpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"order-service/internal/usecase"
)

type OrderHandler struct {
	orderpb.UnimplementedOrderServiceServer
	uc *usecase.OrderUseCase
}

func NewOrderHandler(uc *usecase.OrderUseCase) *OrderHandler {
	return &OrderHandler{uc: uc}
}

func (h *OrderHandler) SubscribeToOrderUpdates(
	req *orderpb.OrderRequest,
	stream orderpb.OrderService_SubscribeToOrderUpdatesServer,
) error {

	orderID := req.OrderId

	for {
		order, err := h.uc.GetByID(fmt.Sprintf("%d", orderID))
		if err != nil {
			return err
		}

		err = stream.Send(&orderpb.OrderStatusUpdate{
			OrderId:   orderID,
			Status:    order.Status,
			UpdatedAt: timestamppb.Now(),
		})
		if err != nil {
			return err
		}

		if order.Status == "Paid" || order.Status == "Failed" {
			return nil
		}

		time.Sleep(2 * time.Second)
	}
}
