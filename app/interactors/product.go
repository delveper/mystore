package interactors

import (
	"context"

	"github.com/delveper/mystore/app/entities"
	"github.com/delveper/mystore/lib/lgr"
)

type ProductInteractor struct {
	repo   ProductRepo
	logger lgr.Logger
}

func NewProductInteractor(repo ProductRepo, logger lgr.Logger) ProductInteractor {
	return ProductInteractor{
		repo:   repo,
		logger: logger,
	}
}

func (p *ProductInteractor) Add(ctx context.Context, prod entities.Product) (int, error) {
	return p.repo.Insert(ctx, prod)
}

func (p *ProductInteractor) Find(ctx context.Context, prod entities.Product) (*entities.Product, error) {
	return p.repo.Select(ctx, prod)
}

func (p *ProductInteractor) FindMany(ctx context.Context) ([]entities.Product, error) {
	return p.repo.SelectMany(ctx)
}

func (p *ProductInteractor) Modify(ctx context.Context, prod entities.Product) error {
	return p.repo.Update(ctx, prod)
}

func (p *ProductInteractor) Remove(ctx context.Context, prod entities.Product) error {
	return p.repo.Delete(ctx, prod)
}
