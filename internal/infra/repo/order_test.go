package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"

	"goffermart/internal/core/model"
)

func TestOrderRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewOrderRepo(db)
	orderID := "123456"

	t.Run("Test CreateOrder with success", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO orders\(id, user_id\) VALUES\(\$1, \$2\)`).
			WithArgs(orderID, 1).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.CreateOrder(context.Background(), orderID, &model.User{ID: 1})
		assert.NoError(t, err)
	})

	t.Run("Test CreateOrder with error", func(t *testing.T) {
		mock.ExpectExec(`INSERT INTO orders\(id, user_id\) VALUES\(\$1, \$2\)`).
			WithArgs("test", 1).
			WillReturnError(errors.New("test_error"))

		err := repo.CreateOrder(context.Background(), "test", &model.User{ID: 1})
		assert.Error(t, err)
	})

	t.Run("Test ListOrders with success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, status, accrual, created_at FROM orders WHERE user_id=\$1 ORDER BY created_at`).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "accrual", "created_at"}).AddRow(orderID, 1, "NEW", 0, time.Now()))

		orders, err := repo.ListOrders(context.Background(), &model.User{ID: 1})
		assert.NoError(t, err)
		assert.NotNil(t, orders)
	})

	t.Run("Test ListOrders with error", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, status, accrual, created_at FROM orders WHERE user_id=\$1 ORDER BY created_at`).
			WithArgs(1).
			WillReturnError(errors.New("test_error"))

		orders, err := repo.ListOrders(context.Background(), &model.User{ID: 1})
		assert.Error(t, err)
		assert.Nil(t, orders)
	})

	t.Run("Test ListOrdersForProcessing with success", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, user_id, status, accrual, created_at FROM orders WHERE status NOT IN \('INVALID', 'PROCESSED'\)`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "status", "accrual", "created_at"}).AddRow("test", 1, "NEW", 0, time.Now()))

		orders, err := repo.ListOrdersForProcessing(context.Background())
		assert.NoError(t, err)
		assert.NotNil(t, orders)
	})
}
