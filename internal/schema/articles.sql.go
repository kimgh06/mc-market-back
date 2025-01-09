// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: articles.sql

package schema

import (
	"context"
	"database/sql"
)

const countArticles = `-- name: CountArticles :one
select count(*)
from articles
`

func (q *Queries) CountArticles(ctx context.Context) (int64, error) {
	row := q.db.QueryRowContext(ctx, countArticles)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createArticle = `-- name: CreateArticle :one
insert into articles (id, title, content, author, head)
values ($1, $2, $3, $4, $5)
returning id, title, content, created_at, updated_at, index, author, head
`

type CreateArticleParams struct {
	ID      int64          `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Author  int64          `json:"author"`
	Head    sql.NullString `json:"head"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (*Article, error) {
	row := q.db.QueryRowContext(ctx, createArticle,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.Author,
		arg.Head,
	)
	var i Article
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Content,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Index,
		&i.Author,
		&i.Head,
	)
	return &i, err
}

const deleteArticle = `-- name: DeleteArticle :exec
delete
from articles
where id = $1
`

func (q *Queries) DeleteArticle(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, deleteArticle, id)
	return err
}

const getArticle = `-- name: GetArticle :one
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author, articles.head, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash
from articles
         left join public.users u on u.id = articles.author
where articles.id = $1
`

type GetArticleRow struct {
	Article Article `json:"article"`
	User    User    `json:"user"`
}

func (q *Queries) GetArticle(ctx context.Context, id int64) (*GetArticleRow, error) {
	row := q.db.QueryRowContext(ctx, getArticle, id)
	var i GetArticleRow
	err := row.Scan(
		&i.Article.ID,
		&i.Article.Title,
		&i.Article.Content,
		&i.Article.CreatedAt,
		&i.Article.UpdatedAt,
		&i.Article.Index,
		&i.Article.Author,
		&i.Article.Head,
		&i.User.ID,
		&i.User.Nickname,
		&i.User.Permissions,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
		&i.User.Cash,
	)
	return &i, err
}

const listArticles = `-- name: ListArticles :many
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author, articles.head, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash
from articles
         left join public.users u on u.id = articles.author
where index > $2::int
order by articles.created_at desc
limit $1
`

type ListArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}

type ListArticlesRow struct {
	Article Article `json:"article"`
	User    User    `json:"user"`
}

func (q *Queries) ListArticles(ctx context.Context, arg ListArticlesParams) ([]*ListArticlesRow, error) {
	rows, err := q.db.QueryContext(ctx, listArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []*ListArticlesRow
	for rows.Next() {
		var i ListArticlesRow
		if err := rows.Scan(
			&i.Article.ID,
			&i.Article.Title,
			&i.Article.Content,
			&i.Article.CreatedAt,
			&i.Article.UpdatedAt,
			&i.Article.Index,
			&i.Article.Author,
			&i.Article.Head,
			&i.User.ID,
			&i.User.Nickname,
			&i.User.Permissions,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.User.Cash,
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
