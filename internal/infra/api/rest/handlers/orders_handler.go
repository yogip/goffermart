package handlers

import (
	"errors"
	"fmt"
	"goffermart/internal/core/service"
	"goffermart/internal/infra/repo"
	"goffermart/internal/logger"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type OrdersHandler struct {
	service *service.OrderService
}

func NewOrdersHandler(service *service.OrderService) *OrdersHandler {
	return &OrdersHandler{service: service}
}

func (h *OrdersHandler) RegisterOrder(ctx *gin.Context) {
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Log.Warn("Error reading body for RegisterOrder route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error reading body: %s", err)},
		)
		return
	}
	orderId, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		logger.Log.Warn("Error reading body for RegisterOrder route", zap.Error(err))
		ctx.JSON(
			http.StatusUnprocessableEntity,
			gin.H{"status": false, "message": fmt.Sprintf("Error converting body: %s", err)},
		)
		return
	}

	log := logger.Log.With(
		zap.Int64("OrderId", orderId),
	)
	log.Debug("Got Order to registration")

	err = h.service.CreateOrder(ctx, orderId)
	logger.Log.Debug(fmt.Sprintf("Order registration result: %v", err))

	// successfull requests
	if err == nil {
		ctx.JSON(
			http.StatusAccepted,
			gin.H{"status": true, "message": "The order has been accepted for processing"},
		)
		return
	}
	if errors.Is(err, repo.ErrOrderAlreadyRegisteredByUser) {
		ctx.JSON(
			http.StatusOK,
			gin.H{"status": true, "message": "The order has already been accepted for processing"},
		)
		return
	}

	// errors
	if errors.Is(err, repo.ErrOrderAlreadyRegisteredByOther) {
		ctx.JSON(
			http.StatusConflict,
			gin.H{"status": false, "message": "The order has been registered by other user"},
		)
		return
	}
	ctx.JSON(
		http.StatusInternalServerError,
		gin.H{"status": false, "message": err},
	)
}

func (h *OrdersHandler) GetOrderList(ctx *gin.Context) {
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
