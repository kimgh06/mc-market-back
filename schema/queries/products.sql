-- name: CreateProduct :one
insert into products (id, creator, name, description, usage, category, price, price_discount, details, tags)
values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
returning *;

-- name: UpdateProduct :one
update products
set creator        = coalesce(sqlc.narg('creator'), creator),
    name           = coalesce(sqlc.narg('name'), name),
    description    = coalesce(sqlc.narg('description'), description),
    usage          = coalesce(sqlc.narg('usage'), usage),
    details        = coalesce(sqlc.narg('details'), details),
    category       = coalesce(sqlc.narg('category'), category),
    tags           = coalesce(sqlc.narg('tags'), tags),
    price          = coalesce(sqlc.narg('price'), price),
    price_discount = coalesce(sqlc.narg('price_discount'), price_discount)
where id = $1
returning *;

-- name: DeleteProduct :exec
delete
from products
where id = $1;

-- name: GetProductById :one
select sqlc.embed(products), sqlc.embed(u), count(pu)
from products
         left join public.users u on u.id = products.creator
         left join public.purchases pu on pu.product = products.id
where products.id = $1
group by products.id, u.id, pu.product;

-- name: ListProducts :many
select sqlc.embed(products), sqlc.embed(u), count(pu)
from products
         left join public.users u on u.id = products.creator
         left join public.purchases pu on pu.product = products.id
where products.id > sqlc.arg('offset')::int
group by products.id, products.created_at, u.id, pu.product
order by products.created_at desc
limit sqlc.arg('limit');