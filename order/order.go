package order

import (
	"github.com/google/uuid"
	"time"
)

type Order struct {
	OrderID     string     `json:"order_id" gorm:"primaryKey"`
	CustomerID  uuid.UUID  `json:"customer_id"`
	LineItems   []LineItem `json:"line_items" gorm:"foreignKey:OrderID"`
	CreatedAt   *time.Time `json:"created_at"`
	ShippedAt   *time.Time `json:"shipped_at"`
	CompletedAt *time.Time `json:"completed_at"`
}

type LineItem struct {
	ItemID   uuid.UUID `json:"item_id" gorm:"primaryKey"`
	Quantity uint      `json:"quantity"`
	Price    uint      `json:"price"`
	OrderID  string    `json:"-"`
}
