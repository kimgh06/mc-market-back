package schema

import (
	"context"
	"fmt"
	"time"
)

type CreateProductUpdateLogParams struct {
	ProductID int64     `json:"product_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type ProductUpdateLog struct {
	ID        int32     `json:"id"`
	ProductID int64    `json:"product_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

const createUpdateLog = `-- name: CreateUpdateLog :one
insert into products_update_log (product_id, title, content, updated_at)
values ($1, $2, $3, now())
returning id, product_id, title, content, updated_at
`

func (q *Queries) CreateUpdateLog(ctx context.Context, arg CreateProductUpdateLogParams) (*ProductUpdateLog, error) {
	row := q.db.QueryRowContext(ctx, createUpdateLog,
		arg.ProductID,
		arg.Title,
		arg.Content,
	)
	var log ProductUpdateLog
	err := row.Scan(
		&log.ID,
		&log.ProductID,
		&log.Title,
		&log.Content,
		&log.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if log.ProductID != arg.ProductID {
		return nil, fmt.Errorf("mismatched product ID: expected %d, got %d", arg.ProductID, log.ProductID)
	}
	return &log, nil
}

type UpdateProductUpdateLogParams struct {
	ID        int32     `json:"id"`
	ProductID int64    `json:"product_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UpdatedAt time.Time `json:"updated_at"`
}

const updateUpdateLog = `-- name: UpdateUpdateLog :one
update products_update_log
set title = $1, content = $2, updated_at = $3
where id = $4
returning id, product_id, title, content, updated_at
`

func (q *Queries) UpdateUpdateLog(ctx context.Context, arg UpdateProductUpdateLogParams) (*ProductUpdateLog, error) {
	row := q.db.QueryRowContext(ctx, updateUpdateLog,
		arg.Title,
		arg.Content,
		arg.UpdatedAt,
		arg.ID,
	)
	var log ProductUpdateLog
	err := row.Scan(
		&log.ID,
		&log.ProductID,
		&log.Title,
		&log.Content,
		&log.UpdatedAt,
	)
	return &log, err
}

const deleteUpdateLog = `-- name: DeleteUpdateLog :exec
delete
from products_update_log
where id = $1
`

func (q *Queries) DeleteUpdateLog(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteUpdateLog, id)
	return err
}

const listUpdateLogsbyProductID = `-- name: ListUpdateLogsByProductID :many
select id, product_id, title, content, updated_at
from products_update_log
where product_id = $1
order by updated_at desc
`

func (q *Queries) ListUpdateLogsByProductID(ctx context.Context, productID int64) ([]ProductUpdateLog, error) {
	rows, err := q.db.QueryContext(ctx, listUpdateLogsbyProductID, productID)
	if err != nil {
		return nil, err
	}
	var logs []ProductUpdateLog
	defer rows.Close()
	for rows.Next() {
		var log ProductUpdateLog
		if err := rows.Scan(
			&log.ID,
			&log.ProductID,
			&log.Title,
			&log.Content,
			&log.UpdatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}
	return logs, nil
}

const getOneUpdateLog = `-- name: GetOneUpdateLog :one
select id, product_id, title, content, updated_at
from products_update_log
where id = $1
`

func (q *Queries) GetOneUpdateLog(ctx context.Context, id int32) (*ProductUpdateLog, error) {
	row := q.db.QueryRowContext(ctx, getOneUpdateLog, id)
	var log ProductUpdateLog
	err := row.Scan(
		&log.ID,
		&log.ProductID,
		&log.Title,
		&log.Content,
		&log.UpdatedAt,
	)
	return &log, err
}

type CheckIsUserOwnerParams struct {
	ProductID int64 `json:"product_id"`
	UserID    int32  `json:"user_id"`
}

const checkIsUserOwner = `-- name: CheckIsUserOwner :one
select (select creator from products where id = $1) = $2
`

func (q *Queries) CheckIsUserOwner(ctx context.Context, arg CheckIsUserOwnerParams) (bool, error) {
	var isOwner bool
	err := q.db.QueryRowContext(ctx, checkIsUserOwner, arg.ProductID, arg.UserID).Scan(&isOwner)
	return isOwner, err
}