package handlers

import (
	"goffermart/internal/core/service"
	"goffermart/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type BalanceHandler struct {
	service *service.IAM
}

func NewBalanceHandler(service *service.IAM) *BalanceHandler {
	return &BalanceHandler{service: service}
}

func (h *BalanceHandler) GetCurrentBalance(ctx *gin.Context) {
	log := logger.Log.With(
		zap.String("tmpl", "---"),
	)

	log.Debug("GetCurrentBalance handler")
	ctx.JSON(
		http.StatusOK,
		gin.H{"current": 500.5, "withdrawn": 42},
	)
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
