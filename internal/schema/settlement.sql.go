package schema

import (
	"context"
)

const createSettlement = `-- name: CreateSettlement :one
insert into settlements (user_id, amount)
values ($1, $2)
returning id, user_id, amount, month, settled
`

type SettlementParams struct {
	UserId int64 `json:"user_id"`
	Amount int32 `json:"amount"`
}

func (q *Queries) CreateSettlement(ctx context.Context, arg SettlementParams) (*Settlement, error) {
	row := q.db.QueryRowContext(ctx, createSettlement, arg.UserId, arg.Amount)
	var i Settlement
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Amount,
		&i.Month,
		&i.Settled,
	)
	return &i, err
}

const increaseSettlement = `-- name: IncreaseSettlement :one
update settlements 
set amount = amount + $2
where user_id = $1
returning id, user_id, amount, month, settled
`

func (q *Queries) IncreaseSettlement(ctx context.Context, arg SettlementParams) (*Settlement, error) {
	row := q.db.QueryRowContext(ctx, increaseSettlement, arg.UserId, arg.Amount)
	var i Settlement
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Amount,
		&i.Month,
		&i.Settled,
	)
	return &i, err
}

const decreaseSettlement = `-- name: IncreaseSettlement :one
update settlements 
set amount = amount - $2
where user_id = $1
returning id, user_id, amount, month, settled
`

func (q *Queries) DecreaseSettlement(ctx context.Context, arg SettlementParams) (*Settlement, error) {
	row := q.db.QueryRowContext(ctx, decreaseSettlement, arg.UserId, arg.Amount)
	var i Settlement
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Amount,
		&i.Month,
		&i.Settled,
	)
	return &i, err
}

const getSettlementByUserIdAndMonth = `-- name: GetSettlementByUserIdAndMonth :one
select id, user_id, amount, month, settled
from settlements
where user_id = $1 and month = $2
`

type GetSettlementParams struct {
	UserId int64  `json:"user_id"`
	Month  string `json:"month"`
}

func (q *Queries) GetSettlementByUserIdAndMonth(ctx context.Context, arg GetSettlementParams) (*Settlement, error) {
	row := q.db.QueryRowContext(ctx, getSettlementByUserIdAndMonth, arg.UserId, arg.Month)
	var i Settlement
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Amount,
		&i.Month,
		&i.Settled,
	)
	return &i, err
}

const getSettlementsByUserIdAndMonth = `-- name: GetSettlementByUserIdAndMonth :one
select id, user_id, amount, month, settled
from settlements
where user_id = $1 AND month >= $2 AND month <= $3
`

type GetSettlementsParams struct {
	UserId     int64  `json:"user_id"`
	StartMonth string `json:"start_month"`
	EndMonth   string `json:"end_month"`
}

func (q *Queries) GetSettlementsByUserIdAndMonths(ctx context.Context, arg GetSettlementsParams) ([]Settlement, error) {
	rows, err := q.db.QueryContext(ctx, getSettlementsByUserIdAndMonth, arg.UserId, arg.StartMonth, arg.EndMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []Settlement
	for rows.Next() {
		var i Settlement
		if err := rows.Scan(
			&i.ID,
			&i.UserId,
			&i.Amount,
			&i.Month,
			&i.Settled,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
}

const settle = `-- name: Settle :one
update settlements
set settled = true
where user_id = $1 and month = $2
returning id, user_id, amount, month, settled`

type SettleParams struct {
	UserId int64  `json:"user_id"`
	Month  string `json:"month"`
}

func (q *Queries) Settle(ctx context.Context, arg SettleParams) (*Settlement, error) {
	row := q.db.QueryRowContext(ctx, settle, arg.UserId, arg.Month)
	var i Settlement
	err := row.Scan(
		&i.ID,
		&i.UserId,
		&i.Amount,
		&i.Month,
		&i.Settled,
	)
	return &i, err
}
