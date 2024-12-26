-- name: CreateArticle :one
insert into articles (id, title, content, author)
values ($1, $2, $3, $4)
returning *;

-- name: GetArticle :one
select sqlc.embed(articles), sqlc.embed(u)
from articles
         left join public.users u on u.id = articles.author
where articles.id = $1;

-- name: ListArticles :many
select sqlc.embed(articles), sqlc.embed(u)
from articles
         left join public.users u on u.id = articles.author
where index > sqlc.arg('offset')::int
order by articles.created_at desc
limit $1;

-- name: DeleteArticle :exec
delete
from articles
where id = $1;

-- name: CountArticles :one
select count(*)
from articles;