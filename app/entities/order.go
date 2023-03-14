package entities

import "time"

type OrderStatus string

const (
	Pending    OrderStatus = "pending"
	Processing OrderStatus = "processing"
	Shipped    OrderStatus = "shipped"
	Delivered  OrderStatus = "delivered"
	Canceled   OrderStatus = "canceled"
)

type Order struct {
	ID         int         `json:"id"`
	CustomerID int         `json:"customer_id"`
	Status     OrderStatus `json:"status"`
	CreatedAt  time.Time   `json:"created_at"`
	DeletedAt  *time.Time  `json:"deleted_at,omitempty"`
}

type OrderItem struct {
	OrderID   int   `json:"order_id"`
	ProductID int   `json:"product_id"`
	Quantity  int   `json:"quantity"`
	Price     int64 `json:"price"`
}
