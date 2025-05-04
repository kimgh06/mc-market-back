package schema

import (
	"context"
	"time"
)

type Adcard struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	LinkURL   string    `json:"link_url"`
	CreatedAt time.Time `json:"created_at"`
	IndexNum  int       `json:"index_num"`
}

const createAdcard = `-- name: CreateAdcard :one
insert into adcard (title, image_url, link_url, index_num, created_at)
values ($1, $2, $3, $4, now())
returning id, title, image_url, link_url, created_at, index_num
`

type CreateAdcardParams struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func (q *Queries) CreateAdcard(ctx context.Context, arg CreateAdcardParams) (*Adcard, error) {
	row := q.db.QueryRowContext(ctx, createAdcard,
		arg.Title,
		arg.ImageURL,
		arg.LinkURL,
		arg.IndexNum,
	)
	var i Adcard
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ImageURL,
		&i.LinkURL,
		&i.CreatedAt,
		&i.IndexNum,
	)
	return &i, err
}

const getAdcard = `-- name: GetAdcard :one
select id, title, image_url, link_url, created_at, index_num
from adcard
where id = $1
`

func (q *Queries) GetAdcard(ctx context.Context, id int64) (*Adcard, error) {
	row := q.db.QueryRowContext(ctx, getAdcard, id)
	var i Adcard
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ImageURL,
		&i.LinkURL,
		&i.CreatedAt,
		&i.IndexNum,
	)
	return &i, err
}

const updateAdcard = `-- name: UpdateAdcard :one
update adcard
set title = $2, image_url = $3, link_url = $4, index_num = $5
where id = $1
returning id, title, image_url, link_url, created_at, index_num
`

type UpdateAdcardParams struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func (q *Queries) UpdateAdcard(ctx context.Context, arg UpdateAdcardParams) (*Adcard, error) {
	row := q.db.QueryRowContext(ctx, updateAdcard,
		arg.ID,
		arg.Title,
		arg.ImageURL,
		arg.LinkURL,
		arg.IndexNum,
	)
	var i Adcard
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.ImageURL,
		&i.LinkURL,
		&i.CreatedAt,
		&i.IndexNum,
	)
	return &i, err
}

const deleteAdcard = `-- name: DeleteAdcard :exec
delete from adcard where id = $1
`

func (q *Queries) DeleteAdcard(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteAdcard, id)
	return err
}

const listAdcards = `-- name: ListAdcards :many
select id, title, image_url, link_url, created_at, index_num
from adcard
order by index_num
`

func (q *Queries) ListAdcards(ctx context.Context) ([]Adcard, error) {
	rows, err := q.db.QueryContext(ctx, listAdcards)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Adcard
	for rows.Next() {
		var i Adcard
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.ImageURL,
			&i.LinkURL,
			&i.CreatedAt,
			&i.IndexNum,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}
