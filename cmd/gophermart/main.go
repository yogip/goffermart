package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"goffermart/internal/core/config"
	"goffermart/internal/core/service"
	"goffermart/internal/infra/api/rest"
	"goffermart/internal/infra/repo"
	"goffermart/internal/logger"
	"goffermart/migrations"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	err = logger.Initialize(cfg.Server.LogLevel)
	if err != nil {
		log.Fatal(err)
	}

	err = migrations.RunMigration(ctx, &cfg.Server)
	if err != nil {
		logger.Log.Fatal("Migration error", zap.String("error", err.Error()))
	}

	if err := run(ctx, cfg); err != nil {
		logger.Log.Fatal("Running server Error", zap.String("error", err.Error()))
	}
}

func run(ctx context.Context, cfg *config.Config) error {
	db, err := sql.Open("pgx", cfg.Server.DatabaseDSN)
	if err != nil {
		return fmt.Errorf("failed to initialize Database: %w", err)
	}
	defer db.Close()

	cancelCtx, cancel := context.WithCancel(ctx)

	repoUser := repo.NewUserRepo(db)
	repoOrders := repo.NewOrderRepo(db)
	repoBalance := repo.NewBalanceRepo(db)

	iamService := service.NewIAMService(repoUser, &cfg.Server)
	ordersService := service.NewOrdersService(repoOrders, &cfg.Server)
	balanceService := service.NewBalanceService(repoBalance, &cfg.Server)
	logger.Log.Info("Service initialized")

	var procWg sync.WaitGroup
	go service.RunProcessing(cancelCtx, &procWg, cfg, repoOrders, repoBalance)

	api := rest.NewAPI(cfg, iamService, ordersService, balanceService)

	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#manually
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := api.Run(cfg.Server.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Log.Info("Runing server error", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	sdCtx, cancelApi := context.WithTimeout(ctx, 5*time.Second)
	defer cancelApi()
	if err := api.Shutdown(sdCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	logger.Log.Info("Waitng for processing goroutines to finish")
	cancel()
	procWg.Wait()

	logger.Log.Info("Server exiting")
	return nil
}
