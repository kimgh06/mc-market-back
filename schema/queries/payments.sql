-- name: CreatePayment :one
insert into payments (id, agent, order_id, amount)
values ($1, $2, $3, $4)
returning *;

-- name: ApprovePayment :one
update payments
set approved = true
where id = $1
returning *;

-- name: DeletePayment :exec
delete
from payments
where id = $1;

-- name: GetPayment :one
select *
from payments
where id = $1;

-- name: GetPaymentByOrderId :one
select *
from payments
where order_id = $1;