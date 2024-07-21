package handlers

import (
	"goffermart/internal/core/service"
	"goffermart/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	BaseHandler
	service *service.BalanceService
}

func NewBalanceHandler(service *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{service: service}
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

	balance, err := h.service.GetBalance(ctx, user)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
	}
	ctx.JSON(http.StatusOK, balance)
}

func (h *BalanceHandler) WithdrawBonuces(ctx *gin.Context) {
	log := logger.Log.With(
		zap.String("tmpl", "---"),
	)

	log.Debug("WithdrawBonuces handler")
	ctx.JSON(
		http.StatusOK,
		gin.H{"order": 123, "sum": 42},
	)
}
