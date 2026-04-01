package http

import (
	"net/http"
	"payment-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	uc *usecase.PaymentUseCase
}

func NewPaymentHandler(uc *usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{uc: uc}
}

func (h *PaymentHandler) Create(c *gin.Context) {
	var req struct {
		OrderID string `json:"order_id"`
		Amount  int64  `json:"amount"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	payment, err := h.uc.Create(req.OrderID, req.Amount)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":         payment.Status,
		"transaction_id": payment.TransactionID,
	})
}

func (h *PaymentHandler) Get(c *gin.Context) {
	orderID := c.Param("order_id")
	payment, err := h.uc.GetByOrderID(orderID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":         payment.Status,
		"transaction_id": payment.TransactionID,
	})
}
