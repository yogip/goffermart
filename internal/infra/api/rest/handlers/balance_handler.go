package handlers

import (
	"errors"
	"fmt"
	"goffermart/internal/core/model"
	"goffermart/internal/core/service"
	"goffermart/internal/infra/repo"
	"goffermart/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	BaseHandler
	balanceService *service.BalanceService
	ordersService  *service.OrderService
}

func NewBalanceHandler(
	balanceService *service.BalanceService,
	ordersService *service.OrderService,
) *BalanceHandler {
	return &BalanceHandler{balanceService: balanceService, ordersService: ordersService}
}

func (h *BalanceHandler) GetCurrentBalance(ctx *gin.Context) {
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
	log.Debug("GetCurrentBalance handler")

	balance, err := h.balanceService.GetBalance(ctx, user)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}
	ctx.JSON(http.StatusOK, balance)
}

func (h *BalanceHandler) WithdrawBonuces(ctx *gin.Context) {
	req := &model.WithdrawRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding schema for WithdrawBonuces route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}
	log := logger.Log.With(
		zap.Int64("OrderID", req.OrderID),
		zap.Float64("WithdrawnBonuses", req.Sum),
	)

	log.Debug("WithdrawBonuces handler")
	user, err := h.getUser(ctx)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}

	// Create order
	err = h.ordersService.CreateOrder(ctx, req.OrderID, user)
	logger.Log.Debug(fmt.Sprintf("Order registration result: %v", err))
	logger.Log.Debug(fmt.Sprintf("1: %t", err != nil))
	logger.Log.Debug(fmt.Sprintf("2: %t", errors.Is(err, repo.ErrOrderAlreadyRegisteredByUser)))
	logger.Log.Debug(fmt.Sprintf("3: %t", err != nil && !errors.Is(err, repo.ErrOrderAlreadyRegisteredByUser)))

	if err != nil && !errors.Is(err, repo.ErrOrderAlreadyRegisteredByUser) {
		if errors.Is(err, repo.ErrOrderAlreadyRegisteredByOther) {
			ctx.JSON(
				http.StatusConflict,
				gin.H{"status": false, "message": "The order has been registered by other user"},
			)
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": false, "message": err})
		return
	}
	logger.Log.Debug("4")

	// Process withdrawl
	err = h.balanceService.ProcessWithdrawl(ctx, user, req)
	logger.Log.Debug(fmt.Sprintf("ProcessWithdrawl result: %v", err))

	if err != nil && errors.Is(err, repo.ErrNoMoney) {
		ctx.JSON(
			http.StatusPaymentRequired,
			gin.H{"status": false, "message": err.Error()},
		)
		return
	}
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err.Error()},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		gin.H{"status": true, "message": "ok"},
	)
}
