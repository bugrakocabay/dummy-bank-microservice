package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		AccountID: RandomString(5),
		UserID:    RandomString(5),
		Currency:  "EUR",
		Balance:   0,
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.AccountID, account.AccountID)
	require.Equal(t, arg.UserID, account.UserID)
	require.Equal(t, arg.Currency, account.Currency)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQueries.GetAccount(context.Background(), account1.AccountID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.UserID, account2.UserID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account1.Balance, account2.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := UpdateAccountParams{
		AccountID: account1.AccountID,
		Balance:   100,
	}
	account2, err := testQueries.UpdateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.UserID, account2.UserID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.Equal(t, account2.Balance, arg.Balance)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)

	err := testQueries.DeleteAccount(context.Background(), account1.AccountID)
	require.NoError(t, err)

	account2, err := testQueries.GetAccount(context.Background(), account1.AccountID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	accounts, err := testQueries.ListAccounts(context.Background())
	require.NoError(t, err)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}
}
