package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"goffermart/internal/core/model"
	"goffermart/internal/infra/accrual"
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

func (r *BalanceRepo) UpdateBalance(
	ctx context.Context,
	tx *sql.Tx,
	userID int64,
	accrual *accrual.Accrual,
) error {
	fun := func() error {
		_, err := tx.ExecContext(
			ctx,
			`INSERT INTO balance(user_id, current, withdrawn) values($1, $2, 0) 
			ON CONFLICT(user_id) 
			DO UPDATE SET current = excluded.current + balance.current`,
			userID, accrual.Accrual,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("update balance error: %w", err)
	}
	return nil
}

func (r *BalanceRepo) ProcessWithdrawl(ctx context.Context, user *model.User, req *model.WithdrawRequest) error {
	var currentBalance float64

	fun := func() error {
		tx, err := r.db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		// Change balance
		row := tx.QueryRowContext(
			ctx,
			`UPDATE balance 
				SET current = current - $1, withdrawn = withdrawn + $1 
			WHERE user_id=$2 
			RETURNING current`,
			req.Sum, user.ID,
		)
		err = row.Scan(&currentBalance)
		if err != nil {
			return fmt.Errorf("update balance error: %w", err)
		}
		if currentBalance < 0 {
			tx.Rollback()
			return ErrNoMoney
		}

		// Add history item
		_, err = tx.ExecContext(
			ctx,
			"INSERT INTO withdrawals(user_id, order_id, sum) values($1, $2, $3)",
			user.ID, req.OrderID, req.Sum,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("insert history error: %w", err)
		}

		err = tx.Commit()
		return err

	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("withdrawals error: %w", err)
	}
	return nil
}

func (r *BalanceRepo) ListWithdrawls(ctx context.Context, user *model.User) (*[]model.Withdrawal, error) {
	withdrawals := []model.Withdrawal{}

	fun := func() error {
		query := "SELECT order_id, sum, processed_at FROM withdrawals WHERE user_id=$1 ORDER BY processed_at"

		rows, err := r.db.QueryContext(ctx, query, user.ID)
		if err != nil {
			return fmt.Errorf("selecting withdrawals error: %w", err)
		}

		for rows.Next() {
			var o model.Withdrawal

			err = rows.Scan(&o.OrderID, &o.Sum, &o.ProcessedAt)
			if err != nil {
				return fmt.Errorf("read withdrawals error: %w", err)
			}

			withdrawals = append(withdrawals, o)
		}

		err = rows.Err()
		if err != nil {
			return fmt.Errorf("reading withdrawals error: %w", err)
		}
		return nil
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	return &withdrawals, nil
}
