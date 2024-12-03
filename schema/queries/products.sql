-- name: CreateProduct :one
insert into products (id, creator, name, description, usage, category, price)
values ($1, $2, $3, $4, $5, $6, $7)
returning *;

-- name: UpdateProduct :one
update products
set creator     = coalesce(sqlc.narg('creator'), creator),
    name        = coalesce(sqlc.narg('name'), name),
    description = coalesce(sqlc.narg('description'), description),
    usage       = coalesce(sqlc.narg('usage'), usage),
    category    = coalesce(sqlc.narg('category'), category),
    price       = coalesce(sqlc.narg('price'), price)
where id = $1
returning *;

-- name: DeleteProduct :exec
delete
from products
where id = $1;

-- name: GetProductById :one
select sqlc.embed(products), sqlc.embed(u)
from products
         left join public.users u on u.id = products.creator
where products.id = $1;

-- name: ListProducts :many
select sqlc.embed(products), sqlc.embed(u)
from products
         left join public.users u on u.id = products.creator
         left join public.downloads d on d.product_id = products.id
where products.id > sqlc.arg('offset')::int
order by products.created_at desc
limit sqlc.arg('limit');