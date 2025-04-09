// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package schema

import (
	"context"
	"database/sql"
	"time"
)

const createUser = `-- name: CreateUser :one
insert into users (id, nickname, created_at, updated_at)
values ($1, $2, $3, $3)
returning id, nickname, permissions, created_at, updated_at, cash
`

type CreateUserParams struct {
	ID        int64          `json:"id"`
	Nickname  string `json:"nickname"`
	CreatedAt time.Time      `json:"created_at"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.ID, arg.Nickname, arg.CreatedAt)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Permissions,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Cash,
	)
	return &i, err
}

const getUserById = `-- name: GetUserById :one
select id, nickname, permissions, created_at, updated_at, cash
from users
where id = $1
`

func (q *Queries) GetUserById(ctx context.Context, id int64) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Permissions,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Cash,
	)
	return &i, err
}

const getUserByNickname = `-- name: GetUserByNickname :one
select id, nickname, permissions, created_at, updated_at, cash
from users
where nickname = $1
`

func (q *Queries) GetUserByNickname(ctx context.Context, nickname string) (*User, error) {
	row := q.db.QueryRowContext(ctx, getUserByNickname, nickname)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Permissions,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Cash,
	)
	return &i, err
}

const listUsers = `-- name: ListUsers :many
select id, nickname, permissions, created_at, updated_at, cash
from users
where users.id > $1::int
order by users.created_at desc
limit $2
`

type ListUsersParams struct {
	Offset int32 `json:"offset"`
	Limit  int32 `json:"limit"`
}

func (q *Queries) ListUsers(ctx context.Context, arg ListUsersParams) ([]*User, error) {
	rows, err := q.db.QueryContext(ctx, listUsers, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*User
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Nickname,
			&i.Permissions,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Cash,
		); err != nil {
			return nil, err
		}
		items = append(items, &i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateUser = `-- name: UpdateUser :one
update users
set nickname    = coalesce($2, nickname),
    permissions = coalesce($3, permissions),
    cash        = coalesce($4, cash),
    updated_at  = now()
where id = $1
returning id, nickname, permissions, created_at, updated_at, cash
`

type UpdateUserParams struct {
	ID          int64          `json:"id"`
	Nickname    sql.NullString `json:"nickname"`
	Permissions sql.NullInt32  `json:"permissions"`
	Cash        sql.NullInt32  `json:"cash"`
}

func (q *Queries) UpdateUser(ctx context.Context, arg UpdateUserParams) (*User, error) {
	row := q.db.QueryRowContext(ctx, updateUser,
		arg.ID,
		arg.Nickname,
		arg.Permissions,
		arg.Cash,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Nickname,
		&i.Permissions,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Cash,
	)
	return &i, err
}


// upload user image
const uploadUserImage = `-- name: UploadUserImage :one
update users
set image_url = $2
where id = $1
`

type UploadUserImageParams struct {
	ID       int64          `json:"id"`
	ImageURL sql.NullString `json:"image_url"`
}	

func (q *Queries) UploadUserImage(ctx context.Context, arg UploadUserImageParams) error {
	_, err := q.db.ExecContext(ctx, uploadUserImage, arg.ID, arg.ImageURL)
	return err
}

// get user image
const getUserImage = `-- name: GetUserImage :one
select image_url
from users
where id = $1
`

func (q *Queries) GetUserImage(ctx context.Context, id int64) (string, error) {
	row := q.db.QueryRowContext(ctx, getUserImage, id)
	var imageURL string
	err := row.Scan(&imageURL)
	return imageURL, err
}

