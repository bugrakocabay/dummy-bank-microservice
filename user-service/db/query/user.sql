-- name: CreateUser :one
INSERT INTO users (user_id, firstname, lastname, password)
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetUser :one
SELECT *
FROM users
WHERE user_id = $1 LIMIT 1;

-- name: UpdateUserPassword :exec
UPDATE users
set password = sqlc.arg(new_password)
WHERE user_id = sqlc.arg(user_id);