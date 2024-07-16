package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	coreErrors "goffermart/internal/core/errors"
	"goffermart/internal/core/model"
	"goffermart/internal/logger"
	"goffermart/internal/retrier"

	"github.com/jackc/pgx/v5/pgconn"
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

func (r *OrderRepo) ListOrders(ctx context.Context) (*[]model.Order, error) {
	user := &model.User{}

	fun := func() error {
		row := r.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email=$1", 123)
		err := row.Scan(&user.ID, &user.Login, &user.PasswordHash)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	return nil, nil
}
