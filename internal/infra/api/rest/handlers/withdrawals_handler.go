package handlers

import (
	"goffermart/internal/core/service"
	"goffermart/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type WithdrawalsHandler struct {
	service *service.IAM
}

func NewWithdrawalsHandler(service *service.IAM) *WithdrawalsHandler {
	return &WithdrawalsHandler{service: service}
}

func (h *WithdrawalsHandler) GetWithdrawalsList(ctx *gin.Context) {
	log := logger.Log.With(
		zap.String("tmpl", "---"),
	)

	log.Debug("GetWithdrawalsList handler")
	//		[
	//	      {
	//	          "order": "2377225624",
	//	          "sum": 500,
	//	          "processed_at": "2020-12-09T16:09:57+03:00"
	//	      }
	//	  ]
}
