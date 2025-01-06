-- name: CreatePurchase :one
insert into purchases (id, purchaser, product, cost)
values ($1, $2, $3, $4)
returning *;

-- name: GetPurchaseById :one
select *
from purchases
where id = $1;

-- name: GetPurchase :one
select *
from purchases
where purchaser = $1
  and product = $2;

-- name: ListProductPurchases :many
select *
from purchases
where product = $1;

-- name: ListUserPurchases :many
select *
from purchases
where purchaser = $1;

-- name: DeletePurchase :exec
delete
from purchases
where id = $1;

-- name: GetUnclaimedRevenuesOfProduct :one
select coalesce(sum(purchases.cost), 0)
from purchases
where product = $1;

-- name: GetUnclaimedRevenuesOfUser :one
select coalesce(sum(purchases.cost), 0), count(purchases)
from purchases
         left join public.products p on p.id = purchases.product
         left join public.users u on u.id = p.creator
where u.id = $1;

-- name: GetUnclaimedPurchasesOfUser :many
select purchases.cost, purchases.purchased_at, sqlc.embed(p)
from purchases
         left join public.products p on p.id = purchases.product
         left join public.users u on u.id = p.creator
where u.id = $1
  and claimed = false;