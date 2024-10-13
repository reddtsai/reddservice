package auth

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestNewService(t *testing.T) {
	db, _, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := NewService(db)
	assert.NotNil(t, svc)
	assert.Equal(t, db, svc.db)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	svc := NewService(db)

	input := CreateUserInput{
		Account: "test_account",
		Email:   "test_account@example.com",
	}

	mock.ExpectPrepare(`INSERT INTO users \("ACCOUNT", "EMAIL", "PASSWORD", "CREATED_AT", "UPDATED_AT"\) VALUES \(\$1, \$2, \$3, \$4, \$5\) RETURNING "ID"`).
		ExpectQuery().
		WithArgs(input.Account, input.Email, USER_DEFAULT_PASSWORD, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"ID"}).AddRow(1))

	userID, err := svc.CreateUser(input)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), userID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
