package service

import (
	"context"
	"errors"
	"goffermart/internal/core/config"
	"goffermart/internal/core/helpers"
	"goffermart/internal/core/model"
	"goffermart/internal/infra/accrual"
	"goffermart/internal/infra/repo"
	"goffermart/internal/logger"

	"time"

	"go.uber.org/zap"
)

func consumer(ctx context.Context, ordersCh chan *model.Order, repo *repo.OrderRepo, client *accrual.AccrualClient, limiter *helpers.StopLimiter) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop consuming orders")
			return
		case order := <-ordersCh:
			limiter.EnsureLimit()
			logger.Log.Debug("Go order for processing", zap.Int64("OrderID", order.ID))
			actual, err := client.GetOrderAccrual(ctx, order.ID)
			var errTM *accrual.ErrorTooManyRequests
			if errors.As(err, &errTM) {
				limiter.StopUntil(errTM.RetryAfter)
			}
			if err != nil {
				logger.Log.Warn("Could not get Accrual information. Skip this order, and continue", zap.Error(err))
				continue
			}

			// order.Status = model.OrderStatus(stat.Status)
			// order.Accrual = actual.Accrual

			err = repo.UpdateAcrual(ctx, actual)
			if err != nil {
				logger.Log.Warn("Could not update Accrual and status. Skip this order, and continue", zap.Error(err))
			}
		}
	}
}

func producer(ctx context.Context, ordersCh chan *model.Order, repo *repo.OrderRepo) {
	orders, err := repo.ListOrdersForProcessing(ctx)
	if err != nil {
		logger.Log.Error("Could not process orders. Skip and try later.", zap.Error(err))
	}
	for _, o := range *orders {
		if ctx.Err() != nil {
			logger.Log.Info("Finish to produce tasks")
			return
		}
		ordersCh <- &o
	}
}

func RunProcessing(ctx context.Context, cfg *config.Config, repo *repo.OrderRepo) {
	logger.Log.Info("Start order processing worker")

	intervalTicker := time.NewTicker(cfg.Accrual.Interval)
	defer intervalTicker.Stop()

	accrual := accrual.NewAccrualClient(cfg)
	limiter := helpers.NewStopLimiter()

	ordersCh := make(chan *model.Order, cfg.Accrual.WorkersCount)
	for i := 0; i < cfg.Accrual.WorkersCount; i++ {
		go consumer(ctx, ordersCh, repo, accrual, limiter)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop processing orders")
			close(ordersCh)
			return
		case <-intervalTicker.C:
			producer(ctx, ordersCh, repo)
		}
	}
}
