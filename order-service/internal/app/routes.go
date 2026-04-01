package app

import (
	orderHttp "order-service/internal/transport/http"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine, handler *orderHttp.OrderHandler) {
	r.POST("/orders", handler.Create)
	r.GET("/orders/:id", handler.Get)
	r.PATCH("/orders/:id/cancel", handler.Cancel)
}
