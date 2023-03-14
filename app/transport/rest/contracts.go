package rest

import (
	"context"
)

type ProductLogic interface {
	Import(context.Context, ) error
	// Fetch(context.Context, models.Book) (models.Book, error)
	// FetchMany(context.Context, models.DataFilter) ([]models.Book, error)
}
