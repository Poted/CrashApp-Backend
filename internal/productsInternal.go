package internal

import (
	"context"
	"go_app/backend/db"
	"go_app/backend/errorz"
	"go_app/backend/models"
)

type IProductsInternal interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
}

type ProductsInternal struct{}

func NewProducts() IProductsInternal {
	return &ProductsInternal{}
}

func (p *ProductsInternal) CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	err := db.Database.
		WithContext(ctx).
		Save(product).
		Error
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return product, nil
}

func (p *ProductsInternal) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	err := db.Database.
		WithContext(ctx).
		Where("id = ?", product.ID).
		Updates(product).
		Error
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return product, nil
}
