package models

import (
	"time"
	"github.com/uptrace/bun"
)

// Order represents a user's purchase order.
type Order struct {
	bun.BaseModel `bun:"table:orders,alias:o"`
	ID        string    `bun:",pk,type:uuid,default:gen_random_uuid()"`
	UserID    string    `bun:"user_id,notnull,type:uuid"`
	Status    string    `bun:"status,notnull"`
	TotalAmount float64 `bun:"total_amount,notnull"`
	CreatedAt  time.Time `bun:",null,default:current_timestamp"`
	UpdatedAt  time.Time `bun:",null,default:current_timestamp"`
}

// OrderItem represents a single item in an order.
type OrderItem struct {
	bun.BaseModel `bun:"table:order_items,alias:oi"`
	ID        string  `bun:",pk,type:uuid,default:gen_random_uuid()"`
	OrderID   string  `bun:"order_id,notnull,type:uuid"`
	ProductID string  `bun:"product_id,notnull,type:uuid"`
	Quantity  int     `bun:"quantity,notnull"`
	UnitPrice float64 `bun:"unit_price,notnull"`
}

