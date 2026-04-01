package app

import (
	paymentHttp "payment-service/internal/transport/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handler *paymentHttp.PaymentHandler) {
	r.POST("/payments", handler.Create)
	r.GET("/payments/:order_id", handler.Get)
}
