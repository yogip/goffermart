package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"goffermart/internal/core/model"
	"goffermart/internal/infra/accrual"
)

func TestBalanceRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewBalanceRepo(db)
	orderID := "123"

	t.Run("Test GetBalance with success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT current, withdrawn FROM balance WHERE user_id=\$1`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"current", "withdrawn"}).AddRow(100.0, 10.0))

		actual, err := repo.GetBalance(context.Background(), &model.User{ID: 1})
		expected := &model.Balance{Current: 100.0, Withdrawn: 10.0}
		assert.NoError(t, err)
		assert.NotNil(t, actual)
		assert.Equal(t, expected, actual)
	})

	t.Run("Test GetBalance with error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT current, withdrawn FROM balance WHERE user_id=\$1`).
			WithArgs(1).
			WillReturnError(errors.New("test_error"))

		balance, err := repo.GetBalance(context.Background(), &model.User{ID: 1})
		assert.Error(t, err)
		assert.Nil(t, balance)
	})

	t.Run("Test UpdateBalance with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO balance\(user_id, current, withdrawn\) values\(\$1, \$2, 0\) 
			ON CONFLICT\(user_id\) 
			DO UPDATE SET current = excluded.current \+ balance.current`).
			WithArgs(1, 100.0).
			WillReturnResult(sqlmock.NewResult(1, 1))

		a := accrual.Accrual{Accrual: 100.0}
		tx, _ := db.BeginTx(context.Background(), nil)
		err := repo.UpdateBalance(context.Background(), tx, 1, &a)
		assert.NoError(t, err)
	})

	t.Run("Test UpdateBalance with error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(`INSERT INTO balance\(user_id, current, withdrawn\) values\(\$1, \$2, 0\) 
			ON CONFLICT\(user_id\) 
			DO UPDATE SET current = excluded.current \+ balance.current`).
			WithArgs(1, 100.0).
			WillReturnError(errors.New("test_error"))

		a := accrual.Accrual{Accrual: 100.0}
		tx, _ := db.BeginTx(context.Background(), nil)
		err := repo.UpdateBalance(context.Background(), tx, 1, &a)
		assert.Error(t, err)
	})

	t.Run("Test ProcessWithdrawl with success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`UPDATE balance 
				SET current = current - \$1, withdrawn = withdrawn \+ \$1 
			WHERE user_id=\$2 
			RETURNING current`).
			WithArgs(100.0, 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"current"}).AddRow(10.0),
			)

		mock.ExpectExec(`INSERT INTO withdrawals\(user_id, order_id, sum\) values\(\$1, \$2, \$3\)`).
			WithArgs(1, orderID, 100.0).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mock.ExpectCommit()

		err := repo.ProcessWithdrawl(context.Background(), &model.User{ID: 1}, &model.WithdrawRequest{Sum: 100.0, OrderID: orderID})
		assert.NoError(t, err)
	})

	t.Run("Test ProcessWithdrawl with ErrNoMoney", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`UPDATE balance
				SET current = current - \$1, withdrawn = withdrawn \+ \$1
			WHERE user_id=\$2
			RETURNING current`).
			WithArgs(100.0, 1).
			WillReturnRows(
				sqlmock.NewRows([]string{"current"}).AddRow(-10.0),
			)
		mock.ExpectRollback()

		err := repo.ProcessWithdrawl(context.Background(), &model.User{ID: 1}, &model.WithdrawRequest{Sum: 100.0, OrderID: "test"})
		assert.ErrorIs(t, err, ErrNoMoney)
	})

	t.Run("Test ProcessWithdrawl with error ", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectQuery(`UPDATE balance
				SET current = current - \$1, withdrawn = withdrawn \+ \$1
			WHERE user_id=\$2
			RETURNING current`).
			WithArgs(100.0, 1).
			WillReturnError(errors.New("test_error"))
		mock.ExpectRollback()

		err := repo.ProcessWithdrawl(context.Background(), &model.User{ID: 1}, &model.WithdrawRequest{Sum: 100.0, OrderID: "test"})
		assert.Error(t, err)
	})

	t.Run("Test ListWithdrawls with success", func(t *testing.T) {
		processed := time.Now()
		mock.ExpectQuery(`SELECT order_id, sum, processed_at FROM withdrawals WHERE user_id=\$1 ORDER BY processed_at`).
			WithArgs(1).
			WillReturnRows(
				sqlmock.NewRows(
					[]string{"order_id", "sum", "processed_at"},
				).AddRow(orderID, 100.0, processed),
			)

		accrual, err := repo.ListWithdrawls(context.Background(), &model.User{ID: 1})
		expected := []model.Withdrawal{
			{
				OrderID:     orderID,
				Sum:         100.0,
				ProcessedAt: processed,
			},
		}
		assert.NoError(t, err)
		assert.NotNil(t, accrual)
		assert.Equal(t, &expected, accrual)
	})

	t.Run("Test ListWithdrawls with error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT order_id, sum, processed_at FROM withdrawals WHERE user_id=\$1 ORDER BY processed_at`).
			WithArgs(1).
			WillReturnError(errors.New("test_error"))

		withdrawals, err := repo.ListWithdrawls(context.Background(), &model.User{ID: 1})
		assert.Error(t, err)
		assert.Nil(t, withdrawals)
	})
}
