package handlers

import (
	"goffermart/internal/core/service"
	"goffermart/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WithdrawalsHandler struct {
	BaseHandler
	service *service.BalanceService
}

func NewWithdrawalsHandler(service *service.BalanceService) *WithdrawalsHandler {
	return &WithdrawalsHandler{service: service}
}

func (h *WithdrawalsHandler) GetWithdrawalList(ctx *gin.Context) {
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
	log.Debug("GetWithdrawalsList handler")

	withdrawals, err := h.service.ListWithdrawls(ctx, user)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": err},
		)
		return
	}
	log.Debug("withdrawals were loaded", zap.Int("Count", len(*withdrawals)))

	if len(*withdrawals) == 0 {
		ctx.String(http.StatusNoContent, "")
		return
	}

	ctx.JSON(
		http.StatusOK,
		withdrawals,
	)
}
