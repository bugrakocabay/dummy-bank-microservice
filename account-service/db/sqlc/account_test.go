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
		Balance:   100,
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

func createRandomTransaction(t *testing.T) Transaction {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	arg := CreateTransactionParams{
		FromAccountID:     account1.AccountID,
		ToAccountID:       account2.AccountID,
		TransactionAmount: 50,
		TransactionID:     RandomString(5),
		Commission:        50 * 0.03,
	}

	transaction, err := testQueries.CreateTransaction(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, transaction)

	require.Equal(t, transaction.FromAccountID, account1.AccountID)
	require.Equal(t, transaction.ToAccountID, account2.AccountID)
	require.Equal(t, transaction.TransactionID, arg.TransactionID)
	require.Equal(t, transaction.TransactionAmount, arg.TransactionAmount)
	require.Equal(t, transaction.Commission, arg.Commission)

	require.NotZero(t, transaction.CreatedAt)

	return transaction
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

func TestAddAccountBalance(t *testing.T) {
	account1 := createRandomAccount(t)

	arg := AddAccountBalanceParams{
		AccountID: account1.AccountID,
		Amount:    100,
	}
	account2, err := testQueries.AddAccountBalance(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account2.Balance, account1.Balance+arg.Amount)
	require.Equal(t, account1.AccountID, account2.AccountID)
	require.Equal(t, account1.UserID, account2.UserID)
	require.Equal(t, account1.Currency, account2.Currency)
	require.WithinDuration(t, account1.CreatedAt, account2.CreatedAt, time.Second)
}

func TestGetAccountBalance(t *testing.T) {
	account := createRandomAccount(t)

	balance, err := testQueries.GetAccountBalance(context.Background(), account.AccountID)
	require.NoError(t, err)
	require.NotEmpty(t, balance)

	require.Equal(t, balance.AccountID, account.AccountID)
	require.Equal(t, balance.Balance, account.Balance)
}

func TestCreateTransaction(t *testing.T) {
	createRandomTransaction(t)
}

func TestGetTransaction(t *testing.T) {
	transaction1 := createRandomTransaction(t)

	transaction2, err := testQueries.GetTransaction(context.Background(), transaction1.TransactionID)
	require.NoError(t, err)
	require.NotEmpty(t, transaction1)

	require.Equal(t, transaction1.TransactionID, transaction2.TransactionID)
	require.Equal(t, transaction1.TransactionAmount, transaction2.TransactionAmount)
	require.Equal(t, transaction1.Commission, transaction2.Commission)
	require.Equal(t, transaction1.ToAccountID, transaction2.ToAccountID)
	require.Equal(t, transaction1.FromAccountID, transaction2.FromAccountID)
	require.WithinDuration(t, transaction1.CreatedAt, transaction2.CreatedAt, time.Second)
}

func TestListTransactions(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomTransaction(t)
	}

	transactions, err := testQueries.ListTransactions(context.Background())
	require.NoError(t, err)

	for _, transaction := range transactions {
		require.NotEmpty(t, transaction)
	}
}
