package schema

import (
	"context"
	"fmt"
	"time"
)

type CreateProductVersionParams struct {
	ProductID   int64  `json:"product_id"`
	VersionName string `json:"version_name"`
	Link        string `json:"link"`
	Index       int    `json:"index"`
}

type ProductVersion struct {
	ID          int32     `json:"id"`
	ProductID   int64     `json:"product_id"`
	VersionName string    `json:"version_name"`
	Link        string    `json:"link"`
	Index       int       `json:"index"`
	UpdatedAt   time.Time `json:"updated_at"`
}

const createProductVersion = `-- name: CreateProductVersion :one
insert into products_versions (product_id, version_name, link, index)
values ($1, $2, $3, $4)
returning id, product_id, version_name, link, index, updated_at
`

func (q *Queries) CreateProductVersion(ctx context.Context, arg CreateProductVersionParams) (*ProductVersion, error) {
	row := q.db.QueryRowContext(ctx, createProductVersion,
		arg.ProductID,
		arg.VersionName,
		arg.Link,
		arg.Index,
	)
	var version ProductVersion
	err := row.Scan(
		&version.ID,
		&version.ProductID,
		&version.VersionName,
		&version.Link,
		&version.Index,
		&version.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if version.ProductID != arg.ProductID {
		return nil, fmt.Errorf("mismatched product ID: expected %d, got %d", arg.ProductID, version.ProductID)
	}
	return &version, nil
}

type UpdateProductVersionParams struct {
	ID          int32  `json:"id"`
	VersionName string `json:"version_name"`
	Link        string `json:"link"`
	Index       int    `json:"index"`
}

const updateProductVersion = `-- name: UpdateProductVersion :one
update products_versions
set version_name = $1, link = $2, index = $3, updated_at = now()
where id = $4
returning id, product_id, version_name, link, index, updated_at
`

func (q *Queries) UpdateProductVersion(ctx context.Context, arg UpdateProductVersionParams) (*ProductVersion, error) {
	row := q.db.QueryRowContext(ctx, updateProductVersion,
		arg.VersionName,
		arg.Link,
		arg.Index,
		arg.ID,
	)
	var version ProductVersion
	err := row.Scan(
		&version.ID,
		&version.ProductID,
		&version.VersionName,
		&version.Link,
		&version.Index,
		&version.UpdatedAt,
	)
	return &version, err
}

const deleteProductVersion = `-- name: DeleteProductVersion :exec
delete
from products_versions
where id = $1
`

func (q *Queries) DeleteProductVersion(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteProductVersion, id)
	return err
}

const listProductVersionsByProductID = `-- name: ListProductVersionsByProductID :many
select id, product_id, version_name, link, index, updated_at
from products_versions
where product_id = $1
order by index
`

func (q *Queries) ListProductVersionsByProductID(ctx context.Context, productID int64) ([]ProductVersion, error) {
	rows, err := q.db.QueryContext(ctx, listProductVersionsByProductID, productID)
	if err != nil {
		return nil, err
	}
	var versions []ProductVersion
	defer rows.Close()
	for rows.Next() {
		var version ProductVersion
		if err := rows.Scan(
			&version.ID,
			&version.ProductID,
			&version.VersionName,
			&version.Link,
			&version.Index,
			&version.UpdatedAt,
		); err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

const getOneProductVersion = `-- name: GetOneProductVersion :one
select id, product_id, version_name, link, index, updated_at
from products_versions
where id = $1
`

func (q *Queries) GetOneProductVersion(ctx context.Context, id int32) (*ProductVersion, error) {
	row := q.db.QueryRowContext(ctx, getOneProductVersion, id)
	var version ProductVersion
	err := row.Scan(
		&version.ID,
		&version.ProductID,
		&version.VersionName,
		&version.Link,
		&version.Index,
		&version.UpdatedAt,
	)
	return &version, err
}
