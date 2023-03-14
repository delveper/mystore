package interactors

import (
	"context"

	"github.com/delveper/mystore/app/entities"
)

type ProductRepo interface {
	Insert(context.Context, entities.Product) error
	Select(context.Context, entities.Product) (*entities.Product, error)
	SelectMany(context.Context) ([]entities.Product, error)
	Update(context.Context, entities.Product) error
	Delete(context.Context, entities.Product) error
}
