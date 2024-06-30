package rest

import (
	"context"
	"goffermart/internal/core/config"
	"goffermart/internal/infra/api/rest/middlewares"
	"goffermart/internal/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type API struct {
	srv *http.Server
}

func NewAPI(cfg *config.Config) *API {
	// serviceHandler := handlers.NewSystemHandler(systemService)

	router := gin.Default()
	router.Use(middlewares.GzipDecompressMiddleware())
	router.Use(middlewares.GzipCompressMiddleware())

	// router.GET("/ping", serviceHandler.Ping)

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
