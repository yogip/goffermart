package rest

import (
	"context"
	"net/http"

	"goffermart/internal/core/config"
	"goffermart/internal/core/service"
	"goffermart/internal/infra/api/rest/handlers"
	"goffermart/internal/infra/api/rest/middlewares"
	"goffermart/internal/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	srv *http.Server
}

func NewAPI(cfg *config.Config, iamService *service.IAM) *API {
	iamHandler := handlers.NewIAMHandler(iamService)

	router := gin.Default()
	router.Use(middlewares.GzipDecompressMiddleware())
	router.Use(middlewares.GzipCompressMiddleware())

	router.POST("/api/user/register", iamHandler.Register)
	router.POST("/api/user/login", iamHandler.Login)

	srv := &http.Server{Handler: router}
	return &API{
		srv: srv,
	}
}

func (api *API) Run(runAddr string) error {
	logger.Log.Info("Run API server", zap.String("Addres", runAddr))
	api.srv.Addr = runAddr
	return api.srv.ListenAndServe()
}

func (api *API) Shutdown(ctx context.Context) error {
	return api.srv.Shutdown(ctx)
}
