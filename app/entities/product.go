package entities

import "time"

type ProductsStatus string

const (
	OutOfStock ProductsStatus = "out_of_stock"
	InStock    ProductsStatus = "in_stock"
	RunningLow ProductsStatus = "running_low"
)

type Product struct {
	ID          int            `json:"id"`
	MerchantID  int            `json:"merchant_id"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Price       int64          `json:"price"`
	Status      ProductsStatus `json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty"`
}
