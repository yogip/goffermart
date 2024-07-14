package handlers

import (
	"goffermart/internal/core/service"
	"goffermart/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	service *service.IAM
}

func NewOrdersHandler(service *service.IAM) *OrdersHandler {
	return &OrdersHandler{service: service}
}

func (h *OrdersHandler) RegisterOrder(ctx *gin.Context) {
	log := logger.Log.With(
		zap.String("tmpl", "---"),
	)

	log.Debug("RegisterOrder handler")
}

func (h *OrdersHandler) GetOrdersList(ctx *gin.Context) {
	log := logger.Log.With(
		zap.String("tmpl", "---"),
	)

	log.Debug("GetOrdersList handler")
	//   [
	// 	{
	// 		"number": "9278923470",
	// 		"status": "PROCESSED",
	// 		"accrual": 500,
	// 		"uploaded_at": "2020-12-10T15:15:45+03:00"
	// 	},
	// 	{
	// 		"number": "12345678903",
	// 		"status": "PROCESSING",
	// 		"uploaded_at": "2020-12-10T15:12:01+03:00"
	// 	},
	// 	{
	// 		"number": "346436439",
	// 		"status": "INVALID",
	// 		"uploaded_at": "2020-12-09T16:09:53+03:00"
	// 	}
	// ]
}
