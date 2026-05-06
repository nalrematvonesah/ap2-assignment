package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"

	pb "github.com/nalrematvonesah/ap2-proto-contracts/gen/go/proto"
)

type PaymentHandler struct {
	client pb.PaymentServiceClient
}

func NewPaymentHandler(conn *grpc.ClientConn) *PaymentHandler {
	return &PaymentHandler{
		client: pb.NewPaymentServiceClient(conn),
	}
}

type PaymentRequest struct {
	OrderID int32   `json:"order_id"`
	Amount  float64 `json:"amount"`
	Email   string  `json:"email"`
}

func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var req PaymentRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := h.client.ProcessPayment(ctx, &pb.PaymentRequest{
		OrderId: req.OrderID,
		Amount:  req.Amount,
		Email:   req.Email,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": res.Status,
	})
}
