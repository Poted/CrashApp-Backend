package internal

import (
	"context"
	"errors"
	"go_app/backend/db"
	"go_app/backend/models"

	"github.com/gofiber/fiber/v2"
)

type IProductsInternal interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, product *models.Product) *fiber.Error
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

func (p *ProductsInternal) DeleteProduct(ctx context.Context, product *models.Product) *fiber.Error {

	query := db.Database.
		WithContext(ctx).
		Delete(product)

	if query.RowsAffected == 0 {
		return fiber.NewError(400, "record not found")
	}

	return fiber.NewError(500, query.Error.Error())
}
