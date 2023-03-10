// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.16.0
// source: user.sql

package db

import (
	"context"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (user_id, firstname, lastname, password)
VALUES ($1, $2, $3, $4) RETURNING id, user_id, firstname, lastname, password, created_at, updated_at
`

type CreateUserParams struct {
	UserID    string `json:"user_id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Password  string `json:"password"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.UserID,
		arg.Firstname,
		arg.Lastname,
		arg.Password,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Firstname,
		&i.Lastname,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, user_id, firstname, lastname, password, created_at, updated_at
FROM users
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, userID string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, userID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Firstname,
		&i.Lastname,
		&i.Password,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateUserPassword = `-- name: UpdateUserPassword :exec
UPDATE users
set password = $1
WHERE user_id = $2
`

type UpdateUserPasswordParams struct {
	NewPassword string `json:"new_password"`
	UserID      string `json:"user_id"`
}

func (q *Queries) UpdateUserPassword(ctx context.Context, arg UpdateUserPasswordParams) error {
	_, err := q.db.ExecContext(ctx, updateUserPassword, arg.NewPassword, arg.UserID)
	return err
}
