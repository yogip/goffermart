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

type UserRepo struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

func NewUserRepo(db *sql.DB) *UserRepo {
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

	repo := &UserRepo{db: db, retrier: ret}

	logger.Log.Info("DB Store initialized")
	return repo
}

func (r *UserRepo) Close() {
	r.db.Close()
}

func (r *UserRepo) CreateUser(ctx context.Context, login string, hashedPassword []byte) (*model.User, error) {
	user := &model.User{Login: login, PasswordHash: nil}
	fun := func() error {
		row := r.db.QueryRowContext(
			ctx,
			"INSERT INTO users(email, password) values($1, $2) RETURNING id",
			login, hashedPassword,
		)
		err := row.Scan(&user.ID)

		// ERROR: duplicate key value violates unique constraint \"users_email_key\" (SQLSTATE 23505)
		var pgErr *pgconn.PgError
		if err != nil && errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrUniqConstrain
		}
		if err != nil {
			return fmt.Errorf("insert to db user: %w", err)
		}
		logger.Log.Debug(fmt.Sprintf("CreateUser %s -> %d", login, user.ID))
		return nil
	}

	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error creating user: %w", err)
	}
	return user, nil
}

func (r *UserRepo) GetUser(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}

	fun := func() error {
		var err error
		row := r.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email=$1", login)
		err = row.Scan(&user.ID, &user.Login, &user.PasswordHash)
		if errors.Is(err, sql.ErrNoRows) {
			return coreErrors.ErrNotFound404
		}
		return err
	}
	err := r.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	return user, nil
}
