-- name: CreateAccount :one
INSERT INTO accounts (account_id, user_id, balance, currency)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetAccount :one
SELECT *
FROM accounts
WHERE account_id = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT *
FROM accounts
WHERE account_id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT *
FROM accounts
ORDER BY id LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
set balance = $2
WHERE id = $1 RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
set balance = balance + sqlc.arg(amount)
WHERE account_id = sqlc.arg(account_id) RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE account_id = $1;

-- name: GetAccountBalance :one
SELECT account_id, balance
FROM accounts
WHERE account_id = $1 LIMIT 1;

-- name: CreateTransaction :one
INSERT INTO transactions (transaction_id, from_account_id, to_account_id, transaction_amount, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetTransaction :one
SELECT *
FROM transactions
WHERE transaction_id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT *
FROM transactions;