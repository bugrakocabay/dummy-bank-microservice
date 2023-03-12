// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: account.sql

package db

import (
	"context"
	"database/sql"
)

const addAccountBalance = `-- name: AddAccountBalance :one
UPDATE accounts
set balance = balance + $1
WHERE account_id = $2 RETURNING id, account_id, user_id, balance, currency, created_at, updated_at
`

type AddAccountBalanceParams struct {
	Amount    int32  `json:"amount"`
	AccountID string `json:"account_id"`
}

func (q *Queries) AddAccountBalance(ctx context.Context, arg AddAccountBalanceParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalance, arg.Amount, arg.AccountID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts (account_id, user_id, balance, currency)
VALUES ($1, $2, $3, $4) RETURNING id, account_id, user_id, balance, currency, created_at, updated_at
`

type CreateAccountParams struct {
	AccountID string `json:"account_id"`
	UserID    string `json:"user_id"`
	Balance   int32  `json:"balance"`
	Currency  string `json:"currency"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount,
		arg.AccountID,
		arg.UserID,
		arg.Balance,
		arg.Currency,
	)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createTransaction = `-- name: CreateTransaction :one
INSERT INTO transactions (transaction_id, from_account_id, to_account_id, transaction_amount, description)
VALUES ($1, $2, $3, $4, $5) RETURNING id, transaction_id, from_account_id, to_account_id, transaction_amount, description, created_at, updated_at
`

type CreateTransactionParams struct {
	TransactionID     string         `json:"transaction_id"`
	FromAccountID     string         `json:"from_account_id"`
	ToAccountID       string         `json:"to_account_id"`
	TransactionAmount int32          `json:"transaction_amount"`
	Description       sql.NullString `json:"description"`
}

func (q *Queries) CreateTransaction(ctx context.Context, arg CreateTransactionParams) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, createTransaction,
		arg.TransactionID,
		arg.FromAccountID,
		arg.ToAccountID,
		arg.TransactionAmount,
		arg.Description,
	)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.TransactionID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.TransactionAmount,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_id = $1
`

func (q *Queries) DeleteAccount(ctx context.Context, accountID string) error {
	_, err := q.db.ExecContext(ctx, deleteAccount, accountID)
	return err
}

const getAccount = `-- name: GetAccount :one
SELECT id, account_id, user_id, balance, currency, created_at, updated_at
FROM accounts
WHERE account_id = $1 LIMIT 1
`

func (q *Queries) GetAccount(ctx context.Context, accountID string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccount, accountID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAccountBalance = `-- name: GetAccountBalance :one
SELECT account_id, balance
FROM accounts
WHERE account_id = $1 LIMIT 1
`

type GetAccountBalanceRow struct {
	AccountID string `json:"account_id"`
	Balance   int32  `json:"balance"`
}

func (q *Queries) GetAccountBalance(ctx context.Context, accountID string) (GetAccountBalanceRow, error) {
	row := q.db.QueryRowContext(ctx, getAccountBalance, accountID)
	var i GetAccountBalanceRow
	err := row.Scan(&i.AccountID, &i.Balance)
	return i, err
}

const getAccountForUpdate = `-- name: GetAccountForUpdate :one
SELECT id, account_id, user_id, balance, currency, created_at, updated_at
FROM accounts
WHERE account_id = $1 LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountForUpdate(ctx context.Context, accountID string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountForUpdate, accountID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getTransaction = `-- name: GetTransaction :one
SELECT id, transaction_id, from_account_id, to_account_id, transaction_amount, description, created_at, updated_at
FROM transactions
WHERE transaction_id = $1 LIMIT 1
`

func (q *Queries) GetTransaction(ctx context.Context, transactionID string) (Transaction, error) {
	row := q.db.QueryRowContext(ctx, getTransaction, transactionID)
	var i Transaction
	err := row.Scan(
		&i.ID,
		&i.TransactionID,
		&i.FromAccountID,
		&i.ToAccountID,
		&i.TransactionAmount,
		&i.Description,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const listAccounts = `-- name: ListAccounts :many
SELECT id, account_id, user_id, balance, currency, created_at, updated_at
FROM accounts
ORDER BY id LIMIT $1
OFFSET $2
`

type ListAccountsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

func (q *Queries) ListAccounts(ctx context.Context, arg ListAccountsParams) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, listAccounts, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Account{}
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.UserID,
			&i.Balance,
			&i.Currency,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listTransactions = `-- name: ListTransactions :many
SELECT id, transaction_id, from_account_id, to_account_id, transaction_amount, description, created_at, updated_at
FROM transactions
`

func (q *Queries) ListTransactions(ctx context.Context) ([]Transaction, error) {
	rows, err := q.db.QueryContext(ctx, listTransactions)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []Transaction{}
	for rows.Next() {
		var i Transaction
		if err := rows.Scan(
			&i.ID,
			&i.TransactionID,
			&i.FromAccountID,
			&i.ToAccountID,
			&i.TransactionAmount,
			&i.Description,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateAccount = `-- name: UpdateAccount :one
UPDATE accounts
set balance = $2
WHERE account_id = $1 RETURNING id, account_id, user_id, balance, currency, created_at, updated_at
`

type UpdateAccountParams struct {
	AccountID string `json:"account_id"`
	Balance   int32  `json:"balance"`
}

func (q *Queries) UpdateAccount(ctx context.Context, arg UpdateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateAccount, arg.AccountID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.UserID,
		&i.Balance,
		&i.Currency,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
