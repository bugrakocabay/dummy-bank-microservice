-- name: GetDailyTransactionReport :one
SELECT COUNT(*) AS num_transactions,
       AVG(transaction_amount) AS avg_transaction_amount,
       SUM(transaction_amount) AS total_transaction_amount,
       SUM(commission) AS total_commission,
       created_at::date AS day
FROM transactions
WHERE created_at::date = $1
GROUP BY day;

-- name: SaveDailyTransactionReport :exec
INSERT INTO daily_transaction_report (num_transactions, avg_transaction_amount, total_transaction_amount, total_commission, day)
VALUES ($1, $2, $3, $4, $5);