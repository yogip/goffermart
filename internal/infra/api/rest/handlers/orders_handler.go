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
	BaseHandler
	service *service.OrderService
}

func NewOrdersHandler(service *service.OrderService) *OrdersHandler {
	return &OrdersHandler{service: service}
}

func (h *OrdersHandler) RegisterOrder(ctx *gin.Context) {
	user, err := h.getUser(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}

	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Log.Warn("Error reading body for RegisterOrder route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error reading body: %s", err)},
		)
		return
	}
	orderID, err := strconv.ParseInt(string(body), 10, 64)
	if err != nil {
		logger.Log.Warn("Error reading body for RegisterOrder route", zap.Error(err))
		ctx.JSON(
			http.StatusUnprocessableEntity,
			gin.H{"status": false, "message": fmt.Sprintf("Error converting body: %s", err)},
		)
		return
	}

	log := logger.Log.With(
		zap.Int64("OrderId", orderID),
		zap.Int64("UserID", user.ID),
		zap.String("Login", user.Login),
	)
	log.Debug("Got Order to registration")

	err = h.service.CreateOrder(ctx, orderID, user)
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
	user, err := h.getUser(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}

	log := logger.Log.With(
		zap.Int64("UserID", user.ID),
		zap.String("Login", user.Login),
	)

	orders, err := h.service.ListOrders(ctx, user)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}
	log.Debug("Orders were loaded", zap.Int("OrdersCount", len(*orders)))

	if len(*orders) == 0 {
		ctx.String(http.StatusNoContent, "")
		return
	}

	ctx.JSON(
		http.StatusOK,
		orders,
	)
}
