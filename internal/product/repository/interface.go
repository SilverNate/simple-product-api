package repository

import (
	"context"
	model "simple-product-api/internal/product"
)

type ProductRepository interface {
	SaveProduct(ctx context.Context, p *model.Product) error
	FindProduct(ctx context.Context, filter model.ListFilter) ([]model.Product, int, error)
	FindProductByID(ctx context.Context, id string) (*model.Product, error)
	FindProductByNameAndType(ctx context.Context, name string, ptype string) (*model.Product, error)
}
