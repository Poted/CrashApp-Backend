package internal

import (
	"context"
	"errors"
	"go_app/backend/db"
	"go_app/backend/models"
)

type IProductsInternal interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProduct(ctx context.Context, product *models.Product) (*models.Product, error)
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
		return nil, err
	}

	return product, nil
}

func (p *ProductsInternal) GetProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	err := db.Database.
		WithContext(ctx).
		First(product).
		Error
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductsInternal) UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error) {

	query := db.Database.
		WithContext(ctx).
		Where("id = ?", &product.ID).
		Updates(product)
	if query.Error != nil {
		return nil, query.Error
	}
	if query.RowsAffected == 0 {
		return nil, errors.New("no rows affected")
	}

	return product, nil
}
