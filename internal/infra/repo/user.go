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

type UserRepo struct {
	db      *sql.DB
	retrier *retrier.Retrier
}

func NewUserRepo(db *sql.DB) *UserRepo {
	// db, err := sql.Open("pgx", cfg.DatabaseDSN)
	// if err != nil {
	// return nil, fmt.Errorf("failed to initialize UserRepo: %w", err)
	// }

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

func (s *UserRepo) Close() {
	s.db.Close()
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
		if err != nil {
			return fmt.Errorf("error creatng user: %w", err)
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

func (s *UserRepo) GetUser(ctx context.Context, login string) (*model.User, error) {
	user := &model.User{}

	fun := func() error {
		var err error
		row := s.db.QueryRowContext(ctx, "SELECT id, email, password FROM users WHERE email=$1", login)
		err = row.Scan(&user.ID, &user.Login)
		if errors.Is(err, sql.ErrNoRows) {
			user = nil
			return nil
		}
		return err
	}
	err := s.retrier.Do(ctx, fun, recoverableErrors...)
	if err != nil {
		return nil, fmt.Errorf("error reading user: %w", err)
	}
	return user, nil
}
