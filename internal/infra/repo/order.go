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

	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
)

type OrderRepo struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

func NewOrderRepo(db *sql.DB) *OrderRepo {
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

	repo := &OrderRepo{db: db, retrier: ret}

	logger.Log.Info("UserRepo initialized")
	return repo
}

func (r *OrderRepo) Tx(ctx context.Context) (*sql.Tx, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (r *OrderRepo) CreateOrder(ctx context.Context, orderId int64, user *model.User) error {
	fun := func() error {
		_, err := r.db.ExecContext(
			ctx,
			"INSERT INTO orders(id, user_id) values($1, $2)",
			orderId, user.ID,
		)

		// ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)
		var pgErr *pgconn.PgError
		if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505" {
			count := 0
			row := r.db.QueryRowContext(ctx, "SELECT count(id) FROM orders WHERE id=$1 AND user_id=$2", orderId, user.ID)
			err := row.Scan(&count)
			if err != nil {
				return fmt.Errorf("create order error: %w", err)
			}
			if count == 0 {
				return ErrOrderAlreadyRegisteredByOther
			}
			return ErrOrderAlreadyRegisteredByUser
		}
		if err != nil {
			return err
		}
		return nil
	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("register order error: %w", err)
	}
	return nil
}

func (r *OrderRepo) UpdateAcrual(ctx context.Context, tx *sql.Tx, accrual *accrual.Accrual) error {
	fun := func() error {
		_, err := tx.ExecContext(
			ctx,
			"UPDATE orders SET status=$2 WHERE id=$",
			accrual.OrderID, accrual.Status,
		)

		if err != nil {
			return err
		}
		return nil
	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return fmt.Errorf("register order error: %w", err)
	}
	return nil
}

func (r *OrderRepo) ListOrders(ctx context.Context, user *model.User) (*[]model.Order, error) {
	orders := []model.Order{}

	fun := func() error {
		query := "SELECT id, user_id, status, accrual, created_at FROM orders WHERE user_id=$1 ORDER BY created_at"

		rows, err := r.db.QueryContext(ctx, query, user.ID)
		if err != nil {
			return fmt.Errorf("selecting orders error: %w", err)
		}

		for rows.Next() {
			var o model.Order

			err = rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Accrual, &o.CreatedAt)
			if err != nil {
				return fmt.Errorf("read order error: %w", err)
			}

			orders = append(orders, o)
		}

		err = rows.Err()
		if err != nil {
			return fmt.Errorf("reading orders error: %w", err)
		}
		return nil
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	return &orders, nil
}

func (r *OrderRepo) ListOrdersForProcessing(ctx context.Context) (*[]model.Order, error) {
	orders := []model.Order{}

	fun := func() error {
		query := "SELECT id, user_id, status, accrual, created_at FROM orders WHERE status NOT IN ('INVALID', 'PROCESSED')"

		rows, err := r.db.QueryContext(ctx, query)
		if err != nil {
			return fmt.Errorf("selecting orders error: %w", err)
		}

		for rows.Next() {
			var o model.Order

			err = rows.Scan(&o.ID, &o.UserID, &o.Status, &o.Accrual, &o.CreatedAt)
			if err != nil {
				return fmt.Errorf("read order error: %w", err)
			}
			logger.Log.Debug("ListOrdersForProcessing", zap.Int64("OrderId", o.ID))
			orders = append(orders, o)
		}

		err = rows.Err()
		if err != nil {
			return fmt.Errorf("reading orders error: %w", err)
		}
		return nil
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}

	return &orders, nil
}
