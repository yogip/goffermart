package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"goffermart/internal/core/model"
	"goffermart/internal/logger"
	"goffermart/internal/retrier"
)

type BalanceRepo struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

func NewBalanceRepo(db *sql.DB) *BalanceRepo {
	ret := &retrier.Retrier{
		Strategy: retrier.Backoff(
			3,             // max attempts
			1*time.Second, // initial delay
			3,             // multiplier
			5*time.Second, // max delay
		),
		OnRetry: func(ctx context.Context, n int, err error) {
			logger.Log.Debug(fmt.Sprintf("Retrying DB. retry #%d: %v", n, err))
		},
	}

	repo := &BalanceRepo{db: db, retrier: ret}

	logger.Log.Info("UserRepo initialized")
	return repo
}

func (r *BalanceRepo) GetBalance(ctx context.Context, user *model.User) (*model.Balance, error) {
	balance := model.Balance{}
	fun := func() error {
		row := r.db.QueryRowContext(ctx, "SELECT current, withdrawn FROM balance WHERE user_id=$1", user.ID)
		err := row.Scan(&balance.Current, &balance.Withdrawn)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return nil
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("get balance error: %w", err)
	}
	return &balance, nil
}
