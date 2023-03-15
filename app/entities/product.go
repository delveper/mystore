package entities

import (
	"errors"
	"fmt"
	"time"
)

type ProductsStatus string

const (
	OutOfStock ProductsStatus = "out_of_stock"
	InStock    ProductsStatus = "in_stock"
	RunningLow ProductsStatus = "running_low"
)

const (
	minTextLength = 2
	maxTextLength = 255
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

// OK checks if Product has valid field values and gather all possible errors.
func (p *Product) OK() error {
	var errList []error

	if p.ID < 0 {
		errList = append(errList, errors.New("id must be a positive integer"))
	}

	if p.MerchantID < 0 {
		errList = append(errList, errors.New("merchant_id must be a positive integer"))
	}

	if len(p.Name) < minTextLength || len(p.Name) > maxTextLength {
		errList = append(errList, fmt.Errorf("name must be between %d and %d characters long", minTextLength, maxTextLength))
	}

	if len(p.Description) > maxTextLength {
		errList = append(errList, fmt.Errorf("description must be less than or equal to %d characters long", maxTextLength))
	}

	if p.Price <= 0 {
		errList = append(errList, errors.New("price must be a positive integer"))
	}

	switch p.Status {
	case OutOfStock, InStock, RunningLow:

	default:
		errList = append(errList, fmt.Errorf("status must be %v | %v | %v", OutOfStock, InStock, RunningLow))
	}

	return errors.Join(errList...)
}
