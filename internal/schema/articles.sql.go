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
insert into articles (id, title, content, author, head, comment_disabled, like_disabled)
values ($1, $2, $3, $4, $5, $6, $7)
returning id, title, content, created_at, updated_at, index, author, head, comment_disabled, like_disabled
`

type CreateArticleParams struct {
	ID      int64          `json:"id"`
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Author  int64          `json:"author"`
	Head    sql.NullString `json:"head"`
	CommentDisabled sql.NullBool `json:"comment_disabled"`
	LikeDisabled sql.NullBool `json:"like_disabled"`
}

func (q *Queries) CreateArticle(ctx context.Context, arg CreateArticleParams) (*Article, error) {
	row := q.db.QueryRowContext(ctx, createArticle,
		arg.ID,
		arg.Title,
		arg.Content,
		arg.Author,
		arg.Head,
		arg.CommentDisabled,
		arg.LikeDisabled,
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
		&i.CommentDisabled,
		&i.LikeDisabled,
	)
	return &i, err
}
const updateArticle = `-- name: UpdateArticle :exec
update articles
set title = $2, content = $3, head = $4, comment_disabled = $5, like_disabled = $6
where id = $1
`

type UpdateArticleParams struct {
	ID              int64          `json:"id"`
	Title           string         `json:"title"`
	Content         string         `json:"content"`
	Head            sql.NullString `json:"head"`
	CommentDisabled bool   `json:"comment_disabled"`
	LikeDisabled    bool   `json:"like_disabled"`
}

func (q *Queries) UpdateArticle(ctx context.Context, arg UpdateArticleParams) error {
	_, err := q.db.ExecContext(ctx, updateArticle, arg.ID, arg.Title, arg.Content, arg.Head, arg.CommentDisabled, arg.LikeDisabled)
	return err
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
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author, 
	(SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) AS head,
	articles.views, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash,
	(select count(*) from articles_likes where articles_likes.article_id = articles.id and articles_likes.kind = true) as likes,
	(select count(*) from articles_likes where articles_likes.article_id = articles.id and articles_likes.kind = false) as disLikes,
	articles.comment_disabled, articles.like_disabled
from articles
		 left join public.users u on u.id = articles.author
where articles.id = $1
`

const updateViews = `-- name: UpdateViews :exec
update articles
set views = views + 1
where id = $1
`

type GetArticleRow struct {
	Article Article `json:"article"`
	User    User    `json:"user"`
	Likes 	int64 	`json:"likes"`
	DisLikes int64 `json:"disLikes"`
	CommentDisabled bool `json:"comment_disabled"`
	LikeDisabled    bool `json:"like_disabled"`
}

func (q *Queries) GetArticle(ctx context.Context, id int64) (*GetArticleRow, error) {
	_, err := q.db.ExecContext(ctx, updateViews, id)
	if err != nil {
		return nil, err
	}
	row := q.db.QueryRowContext(ctx, getArticle, id)
	var i GetArticleRow
	err = row.Scan(
		&i.Article.ID,
		&i.Article.Title,
		&i.Article.Content,
		&i.Article.CreatedAt,
		&i.Article.UpdatedAt,
		&i.Article.Index,
		&i.Article.Author,
		&i.Article.Head,
		&i.Article.Views,
		&i.User.ID,
		&i.User.Nickname,
		&i.User.Permissions,
		&i.User.CreatedAt,
		&i.User.UpdatedAt,
		&i.User.Cash,
		&i.Likes,
		&i.DisLikes,
		&i.CommentDisabled,
		&i.LikeDisabled,
	)
	if err != nil {
		return nil, err
	}
	return &i, err
}

const listNoticeheadArticles = `-- name: ListNotices :many
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author,
	(SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) AS head, 
	articles.views, CASE WHEN articles.content LIKE '%<img src%' THEN TRUE ELSE FALSE END AS has_img, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash, 
	(select count(*) from comments where comments.article_id = articles.id) as comment_count,
	(select count(*) from articles_likes where articles_likes.article_id = articles.id and articles_likes.kind = true) as likes
from articles
		left join public.users u on u.id = articles.author
where (SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) like '공지'
order by articles.created_at desc
`

const listArticles = `-- name: ListArticles :many
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author,
	(SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) AS head, 
	articles.views, CASE WHEN articles.content LIKE '%<img src%' THEN TRUE ELSE FALSE END AS has_img, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash, 
	(select count(*) from comments where comments.article_id = articles.id) as comment_count,
	(select count(*) from articles_likes where articles_likes.article_id = articles.id and articles_likes.kind = true) as likes
from articles
		left join public.users u on u.id = articles.author
where articles.index > $2::int and ((SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) not like '공지' or (SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) is null)
order by articles.created_at desc
limit $1
`

type ListArticlesParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	ArticleID int64 `json:"article_id"`
}

type ListArticlesRow struct {
	Article Article `json:"article"`
	HasImg bool    `json:"has_img"`
	User    User    `json:"user"`
	Likes  int64   `json:"likes"`
	CommentCount int64 `json:"comment_count"`
}

func (q *Queries) ListArticles(ctx context.Context, arg ListArticlesParams) ([]*ListArticlesRow, error) {
	rows, err := q.db.QueryContext(ctx, listNoticeheadArticles)
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
			&i.Article.Views,
			&i.HasImg,
			&i.User.ID,
			&i.User.Nickname,
			&i.User.Permissions,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.User.Cash,
			&i.CommentCount,
			&i.Likes,
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

	rows, err = q.db.QueryContext(ctx, listArticles, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
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
			&i.Article.Views,
			&i.HasImg,
			&i.User.ID,
			&i.User.Nickname,
			&i.User.Permissions,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.User.Cash,
			&i.CommentCount,
			&i.Likes,
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


const listArticlesByHead = `-- name: ListArticlesByHead :many
select articles.id, articles.title, articles.content, articles.created_at, articles.updated_at, articles.index, articles.author,
	(SELECT name FROM article_head_type WHERE article_head_type.id = articles.head::integer) AS head,
	articles.views, CASE WHEN articles.content LIKE '%<img src%' THEN TRUE ELSE FALSE END AS has_img, u.id, u.nickname, u.permissions, u.created_at, u.updated_at, u.cash,
	(select count(*) from comments where comments.article_id = articles.id) as comment_count,
	(select count(*) from articles_likes where articles_likes.article_id = articles.id and articles_likes.kind = true) as likes
from articles
		left join public.users u on u.id = articles.author
where articles.head = $2::text and articles.index > $3::int
order by articles.created_at desc
limit $1
`

type ListArticlesByHeadParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
	Head   string `json:"head"`
}

func (q *Queries) ListArticlesByHead(ctx context.Context, arg ListArticlesByHeadParams) ([]*ListArticlesRow, error) {
	rows, err := q.db.QueryContext(ctx, listArticlesByHead, arg.Limit, arg.Head, arg.Offset)
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
			&i.Article.Views,
			&i.HasImg,
			&i.User.ID,
			&i.User.Nickname,
			&i.User.Permissions,
			&i.User.CreatedAt,
			&i.User.UpdatedAt,
			&i.User.Cash,
			&i.CommentCount,
			&i.Likes,
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

const getArticleAuthor = `-- name: GetArticleAuthor :one
select author
from articles
where id = $1
`

func (q *Queries) GetArticleAuthor(ctx context.Context, id int64) int64 {
	var author int64
	q.db.QueryRowContext(ctx, getArticleAuthor, id).Scan(&author)
	return author
}
