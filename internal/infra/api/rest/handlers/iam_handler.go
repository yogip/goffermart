package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"goffermart/internal/core/model"
	"goffermart/internal/core/service"
	"goffermart/internal/infra/repo"
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
	user := &model.UserRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		logger.Log.Error("Error binding schema for Login route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}
	log := logger.Log.With(
		zap.String("email", user.Login),
	)

	token, err := h.service.Login("token")
	if err != nil {
		log.Warn("Could not Auth user", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}

	log.Debug("User was Authentificated")
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	ctx.JSON(
		http.StatusOK,
		gin.H{"status": false, "message": "Ok"},
	)
}

func (h *IAMHandler) Register(ctx *gin.Context) {
	user := &model.UserRequest{}
	if err := ctx.ShouldBindBodyWithJSON(&user); err != nil {
		logger.Log.Error("Error binding schema for Register route", zap.Error(err))
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"status": false, "message": fmt.Sprintf("Error binding body: %s", err)},
		)
		return
	}
	log := logger.Log.With(
		zap.String("login", user.Login),
	)

	token, err := h.service.Register(ctx, user)
	if err != nil && errors.Is(err, repo.ErrUniqConstrain) {
		log.Warn("Could not register user. User already exists", zap.Error(err))
		ctx.JSON(
			http.StatusConflict,
			gin.H{"status": false, "message": "User already exists"},
		)
		return
	}
	if err != nil {
		log.Warn("Could not register user", zap.Error(err))
		ctx.JSON(
			http.StatusInternalServerError,
			gin.H{"status": false, "message": fmt.Sprintf("User registration error: %s", err)},
		)
		return
	}

	log.Debug("User was register")
	ctx.Header("Authorization", fmt.Sprintf("Bearer %s", token))
	ctx.JSON(
		http.StatusOK,
		gin.H{"status": false, "message": "User was registered"},
	)
}
