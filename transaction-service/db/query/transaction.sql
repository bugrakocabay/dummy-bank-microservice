-- name: CreateTransaction :one
INSERT INTO transactions (transaction_id, from_account_id, to_account_id, transaction_amount, description)
VALUES ($1, $2, $3, $4, $5) RETURNING *;

-- name: GetTransaction :one
SELECT *
FROM transactions
WHERE id = $1 LIMIT 1;

-- name: ListTransactions :many
SELECT *
FROM transactions;