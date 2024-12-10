-- name: CreateUser :one
insert into users (id, nickname, created_at, updated_at)
values ($1, $2, $3, $3)
returning *;

-- name: GetUserById :one
select *
from users
where id = $1;

-- name: GetUserByNickname :one
select *
from users
where nickname = $1;

-- name: ListUsers :many
select *
from users
where users.id > sqlc.arg('offset')::int
order by users.created_at desc
limit sqlc.arg('limit');

-- name: UpdateUser :one
update users
set nickname    = coalesce(sqlc.narg('nickname'), nickname),
    permissions = coalesce(sqlc.narg('permissions'), permissions),
    cash        = coalesce(sqlc.narg('cash'), cash),
    updated_at  = now()
where id = $1
returning *;