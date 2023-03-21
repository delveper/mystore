package interactors

import (
	"context"
	"fmt"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/lib/lgr"
)

// ProductInteractor implements use cases for interacting with entities.Product.
type ProductInteractor struct {
	repo   ProductRepo // implemented in repository layer
	logger lgr.Logger
}

// NewProductInteractor is a constructor function for creating a new ProductInteractor instance.
func NewProductInteractor(repo ProductRepo, logger lgr.Logger) ProductInteractor {
	return ProductInteractor{
		repo:   repo,
		logger: logger,
	}
}

// Add inserts a new entities.Product into the repository and returns its ID.
func (p ProductInteractor) Add(ctx context.Context, prod entities.Product) (int, error) {
	id, err := p.repo.Insert(ctx, prod)
	if err != nil {
		return 0, fmt.Errorf("inserting product into repository: %w", err)
	}

	return id, nil
}

// Find retrieves an entities.Product from the repository based on its ID.
func (p ProductInteractor) Find(ctx context.Context, prod entities.Product) (*entities.Product, error) {
	found, err := p.repo.Select(ctx, prod)
	if err != nil {
		return nil, fmt.Errorf("finding product in repository: %w", err)
	}

	return found, nil
}

// FindMany retrieves all entities.Product from the repository.
func (p ProductInteractor) FindMany(ctx context.Context) ([]entities.Product, error) {
	prods, err := p.repo.SelectMany(ctx)
	if err != nil {
		return nil, fmt.Errorf("retrievieng products from repository: %w", err)
	}

	return prods, nil
}

// Modify updates an existing entities.Product in the repository.
func (p ProductInteractor) Modify(ctx context.Context, prod entities.Product) error {
	err := p.repo.Update(ctx, prod)
	if err != nil {
		return fmt.Errorf("updating product in repository: %w", err)
	}

	return nil
}

// Remove deletes an existing entities.Product entity from the repository.
func (p ProductInteractor) Remove(ctx context.Context, prod entities.Product) (int, error) {
	id, err := p.repo.Delete(ctx, prod)
	if err != nil {
		return 0, fmt.Errorf("deleting product from repository: %w", err)
	}

	return id, nil
}
