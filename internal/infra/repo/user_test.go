package repo

import (
	"context"
	"database/sql"
	"errors"
	coreErrors "goffermart/internal/core/errors"
	"goffermart/internal/core/model"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUserRepo(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	repo := NewUserRepo(db)

	t.Run("Test GetUser with existing user", func(t *testing.T) {
		mock.ExpectQuery(`SELECT id, email, password FROM users WHERE email=\$1`).
			WithArgs("test").
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "email", "password"}).
					AddRow(1, "test", "test_pwd_hash"),
			)

		actual, err := repo.GetUser(context.Background(), "test")
		assert.NoError(t, err)
		assert.NotEmpty(t, actual)

		h := []byte("test_pwd_hash")
		expected := model.User{ID: 1, Login: "test", PasswordHash: &h}
		assert.Equal(t, &expected, actual)
	})

	t.Run("Test GetUser with non-existing user", func(t *testing.T) {
		mock.ExpectQuery("SELECT id, email, password FROM users WHERE email=??").
			WithArgs("test").
			WillReturnError(sql.ErrNoRows)

		user, err := repo.GetUser(context.Background(), "test")
		assert.ErrorIs(t, err, coreErrors.ErrNotFound404)
		assert.Nil(t, user)
	})

	t.Run("Test CreateUser with success", func(t *testing.T) {
		pwdHash := []byte("test_password")
		mock.ExpectQuery(`INSERT INTO users\(email, password\) values\(\$1, \$2\) RETURNING id`).
			WithArgs("test", pwdHash).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

		actual, err := repo.CreateUser(context.Background(), "test", pwdHash)
		expected := model.User{ID: 1, Login: "test", PasswordHash: nil}

		assert.NoError(t, err)
		assert.NotEmpty(t, actual)

		assert.Equal(t, &expected, actual)
	})

	t.Run("Test CreateUser with error", func(t *testing.T) {
		pwdHash := []byte("test_password")
		mock.ExpectExec(`INSERT INTO users\(email, password\) values\(\$1, \$2\) RETURNING id`).
			WithArgs("test", pwdHash).
			WillReturnError(errors.New("test_error"))

		user, err := repo.CreateUser(context.Background(), "test", []byte("test_password"))
		assert.Error(t, err)
		assert.Nil(t, user)
	})
}
