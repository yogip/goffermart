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

func consumer(
	ctx context.Context,
	ordersCh chan model.Order,
	orderRepo *repo.OrderRepo,
	balanceRepo *repo.BalanceRepo,
	client *accrual.AccrualClient,
	limiter *helpers.StopLimiter,
) {
	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop consuming orders")
			return
		case order := <-ordersCh:
			limiter.EnsureLimit()

			log := logger.Log.With(
				zap.String("OrderID", order.ID),
				zap.Int64("UserID", order.UserID),
				zap.String("StatusOld", string(order.Status)),
			)
			log.Debug("Got order for processing")

			acrual, err := client.GetOrderAccrual(ctx, order.ID)
			log.Debug("Got Acrual resp")

			var errTM *accrual.ErrorTooManyRequests
			if errors.As(err, &errTM) {
				limiter.StopUntil(errTM.RetryAfter)
			}
			if err != nil {
				log.Warn("Could not get Accrual information. Skip this order and continue", zap.Error(err))
				continue
			}
			if acrual == nil {
				log.Warn("There is no Accrual information. Skip this order and continue")
				continue
			}
			log = logger.Log.With(
				zap.Float64("Accrual", acrual.Accrual),
				zap.String("StatusNew", string(acrual.Status)),
			)

			// update order status and sum
			tx, err := orderRepo.Tx(ctx)
			if err != nil {
				log.Warn("Could not start transaction. Skip this order and continue", zap.Error(err))
				continue
			}
			err = orderRepo.UpdateAcrual(ctx, tx, acrual)
			if err != nil {
				tx.Rollback()
				log.Warn("Could not UpdateAcrual. Skip this order, and continue", zap.Error(err))
			}

			// increse balance
			err = balanceRepo.UpdateBalance(ctx, tx, order.UserID, acrual)
			if err != nil {
				tx.Rollback()
				log.Warn("Could not UpdateBalance. Skip this order and continue", zap.Error(err))
			}

			tx.Commit()
			// finish order processing
		}
	}
}

func producer(ctx context.Context, ordersCh chan model.Order, repo *repo.OrderRepo) {
	orders, err := repo.ListOrdersForProcessing(ctx)
	if err != nil {
		logger.Log.Error("Could not process orders. Skip and try later.", zap.Error(err))
	}
	for _, o := range *orders {
		if ctx.Err() != nil {
			logger.Log.Info("Finish to produce tasks")
			return
		}
		logger.Log.Debug("Add order to processing queue", zap.String("OrderId", o.ID))
		ordersCh <- o
	}
}

func RunProcessing(
	ctx context.Context,
	cfg *config.Config,
	orderRepo *repo.OrderRepo,
	balanceRepo *repo.BalanceRepo,
) {
	logger.Log.Info("Start order processing worker")

	intervalTicker := time.NewTicker(cfg.Accrual.Interval)
	defer intervalTicker.Stop()

	accrual := accrual.NewAccrualClient(cfg)
	limiter := helpers.NewStopLimiter()

	ordersCh := make(chan model.Order, cfg.Accrual.WorkersCount)
	for i := 0; i < cfg.Accrual.WorkersCount; i++ {
		go consumer(ctx, ordersCh, orderRepo, balanceRepo, accrual, limiter)
	}

	for {
		select {
		case <-ctx.Done():
			logger.Log.Info("Stop processing orders")
			close(ordersCh)
			return
		case <-intervalTicker.C:
			producer(ctx, ordersCh, orderRepo)
		}
	}
}
