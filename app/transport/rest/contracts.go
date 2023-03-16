package rest

import (
	"context"

	"github.com/delveper/mystore/app/entities"
)

//go:generate mockgen -source=contracts.go -destination=mock_product_logic_test.go -package=rest
type ProductLogic interface {
	Add(context.Context, entities.Product) (id int, err error)
	Find(context.Context, entities.Product) (*entities.Product, error)
	FindMany(context.Context) ([]entities.Product, error)
	Modify(context.Context, entities.Product) error
	Remove(context.Context, entities.Product) (id int, err error)
}
