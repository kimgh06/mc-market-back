-- name: CreateProduct :one
insert into products (id, creator, name, description, usage, category)
values ($1, $2, $3, $4, $5, $6)
returning *;

-- name: DeleteProduct :exec
delete
from products
where id = $1;

-- name: GetProductById :one
select *
from products
where id = $1;

-- name: ListProducts :many
select sqlc.embed(products), sqlc.embed(u)
from products
         left join public.users u on u.id = products.creator
where products.id > sqlc.arg('offset')::int
order by products.created_at desc
limit sqlc.arg('limit');