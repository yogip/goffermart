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
	balanceHandler := handlers.NewBalanceHandler(iamService)
	ordersHandler := handlers.NewOrdersHandler(iamService)
	withdrawalsHandler := handlers.NewWithdrawalsHandler(iamService)

	router := gin.Default()
	router.Use(middlewares.GzipDecompressMiddleware())
	router.Use(middlewares.GzipCompressMiddleware())

	router.POST("/api/user/register", iamHandler.Register)
	router.POST("/api/user/login", iamHandler.Login)

	authMiddleware := middlewares.NewAuthMiddleware(iamService)

	// Route group - /api/user/orders
	ordersRoute := router.Group("/api/user/orders", authMiddleware.AuthRequired())
	{
		ordersRoute.POST("/", ordersHandler.RegisterOrder) // register an order
		ordersRoute.GET("/", ordersHandler.GetOrdersList)  // orders list
	}

	// Route group - /api/user/balance
	balanceRoute := router.Group("/api/user/balance", authMiddleware.AuthRequired())
	{
		balanceRoute.GET("/", balanceHandler.GetCurrentBalance)        // get current balance
		balanceRoute.POST("/withdraw", balanceHandler.WithdrawBonuces) // deduct bonuses from your balance
	}

	// Route group - /api/user/withdrawals
	withdrawalsRoute := router.Group("/api/user/withdrawals", authMiddleware.AuthRequired())
	{
		withdrawalsRoute.GET("/", withdrawalsHandler.GetWithdrawalsList) // withdrawals list
	}

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
