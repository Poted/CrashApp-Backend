package internal

import (
	"context"
	"go_app/backend/db"
	"go_app/backend/models"
	"os"

	"github.com/gofiber/fiber/v2"
)

type IProductsInternal interface {
	CreateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	GetProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	UpdateProduct(ctx context.Context, product *models.Product) (*models.Product, error)
	DeleteProduct(ctx context.Context, product *models.Product, withFiles bool) ([]string, *fiber.Error)
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
		return nil, fiber.NewError(400, "record not found")
	}

	return product, nil
}

func (p *ProductsInternal) DeleteProduct(ctx context.Context, product *models.Product, withFiles bool) ([]string, *fiber.Error) {

	tx := db.Database.Begin()

	// productFiles := []string{}

	query := tx.
		WithContext(ctx).
		Delete(product)

	if query.RowsAffected == 0 {
		return nil, fiber.NewError(400, "record not found")
	}

	var deletedFiles []string
	for _, v := range product.Files {
		if err := os.Remove(v.FilePath(true)); err != nil {
			return deletedFiles, fiber.NewError(500, query.Error.Error())
		}
		deletedFiles = append(deletedFiles, v.ID.String())
	}

	// tx.
	// Model(models.ProductFile{}).
	// Where("product_id = ?", product.ID).
	// Select("file_id").
	// Scan(&productFiles)

	err := tx.Commit().Error
	if err != nil {
		return deletedFiles, fiber.NewError(500, query.Error.Error())
	}

	return deletedFiles, nil
}
