package handlers

import (
	"fmt"
	"net/http"

	"goffermart/internal/core/model"
	"goffermart/internal/core/service"
	"goffermart/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type IAMHandler struct {
	service *service.IAM
}

func NewIAMHandler(service *service.IAM) *IAMHandler {
	return &IAMHandler{service: service}
}

func (h *IAMHandler) Login(ctx *gin.Context) {
	token, err := h.service.Login("token")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "OK") // todo
		return
	}

	// Authorization: Bearer <token>
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "OK") // todo
}

func (h *IAMHandler) Register(ctx *gin.Context) {
	req := &model.UserRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
		logger.Log.Error("Error binding schema for Register route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest, 
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)}
		)
		return
	}
	log := logger.Log.With(
		zap.String("name", req.ID),
		zap.String("type", req.MType.String()),
		zap.Int64p("delta", req.Delta),
		zap.Float64p("value", req.Value),
	)

	token, err := h.service.Register("token")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "OK") // todo
		return
	}

	// Authorization: Bearer <token>
	ctx.Header("Authorization", token)
	ctx.String(http.StatusOK, "OK") // todo
}
