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