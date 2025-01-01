// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0

package schema

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type Article struct {
	ID        int64         `json:"id"`
	Title     string        `json:"title"`
	Content   string        `json:"content"`
	CreatedAt time.Time     `json:"created_at"`
	UpdatedAt time.Time     `json:"updated_at"`
	Index     sql.NullInt32 `json:"index"`
	Author    int64         `json:"author"`
}

type Download struct {
	UserID    int64     `json:"user_id"`
	ProductID int64     `json:"product_id"`
	Time      time.Time `json:"time"`
}

type Like struct {
	UserID    sql.NullInt64 `json:"user_id"`
	ProductID sql.NullInt64 `json:"product_id"`
}

type Payment struct {
	ID        int64        `json:"id"`
	Agent     int64        `json:"agent"`
	OrderID   uuid.UUID    `json:"order_id"`
	Amount    int32        `json:"amount"`
	Approved  bool         `json:"approved"`
	CreatedAt time.Time    `json:"created_at"`
	Failed    sql.NullBool `json:"failed"`
}

type Product struct {
	ID            int64         `json:"id"`
	Creator       int64         `json:"creator"`
	Category      string        `json:"category"`
	Name          string        `json:"name"`
	Description   string        `json:"description"`
	Usage         string        `json:"usage"`
	Price         int32         `json:"price"`
	PriceDiscount sql.NullInt32 `json:"price_discount"`
	Ts            interface{}   `json:"ts"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
	Details       string        `json:"details"`
	Tags          []string      `json:"tags"`
}

type Purchase struct {
	ID          int64     `json:"id"`
	Purchaser   int64     `json:"purchaser"`
	Product     int64     `json:"product"`
	PurchasedAt time.Time `json:"purchased_at"`
	Claimed     bool      `json:"claimed"`
	Cost        int32     `json:"cost"`
}

type User struct {
	ID          int64          `json:"id"`
	Nickname    sql.NullString `json:"nickname"`
	Permissions int32          `json:"permissions"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	Cash        int32          `json:"cash"`
}
