package schema

import (
	"context"
	"time"
)

type Banner struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	ImageURL  string    `json:"image_url"`
	LinkURL   string    `json:"link_url"`
	CreatedAt time.Time `json:"created_at"`
	IndexNum  int       `json:"index_num"`
}

const createBanner = `-- name: CreateBanner :one
insert into banner (title, image_url, link_url, index_num, created_at)
values ($1, $2, $3, $4, now())
returning id, title, image_url, link_url, created_at, index_num
`

type CreateBannerParams struct {
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func (q *Queries) CreateBanner(ctx context.Context, arg CreateBannerParams) (*Banner, error) {
	row := q.db.QueryRowContext(ctx, createBanner,
		arg.Title,
		arg.ImageURL,
		arg.LinkURL,
		arg.IndexNum,
	)
	var i Banner
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

const getBanner = `-- name: GetBanner :one
select id, title, image_url, link_url, created_at, index_num
from banner
where id = $1
`

func (q *Queries) GetBanner(ctx context.Context, id int64) (*Banner, error) {
	row := q.db.QueryRowContext(ctx, getBanner, id)
	var i Banner
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

const updateBanner = `-- name: UpdateBanner :one
update banner
set title = $2, image_url = $3, link_url = $4, index_num = $5
where id = $1
returning id, title, image_url, link_url, created_at, index_num
`

type UpdateBannerParams struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	ImageURL string `json:"image_url"`
	LinkURL  string `json:"link_url"`
	IndexNum int    `json:"index_num"`
}

func (q *Queries) UpdateBanner(ctx context.Context, arg UpdateBannerParams) (*Banner, error) {
	row := q.db.QueryRowContext(ctx, updateBanner,
		arg.ID,
		arg.Title,
		arg.ImageURL,
		arg.LinkURL,
		arg.IndexNum,
	)
	var i Banner
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

const deleteBanner = `-- name: DeleteBanner :exec
delete from banner where id = $1
`

func (q *Queries) DeleteBanner(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteBanner, id)
	return err
}

const listBanners = `-- name: ListBanners :many
select id, title, image_url, link_url, created_at, index_num
from banner
order by index_num
`

func (q *Queries) ListBanners(ctx context.Context) ([]Banner, error) {
	rows, err := q.db.QueryContext(ctx, listBanners)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Banner
	for rows.Next() {
		var i Banner
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
