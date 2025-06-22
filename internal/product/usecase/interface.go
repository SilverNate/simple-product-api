package usecase

import (
	"context"
	model "simple-product-api/internal/product"
)

type ProductUsecase interface {
	CreateProduct(ctx context.Context, p *model.Product) error
	ListProduct(ctx context.Context, filter model.ListFilter) ([]model.Product, int, error)
	GetProductByID(ctx context.Context, id string) (*model.Product, error)
}
