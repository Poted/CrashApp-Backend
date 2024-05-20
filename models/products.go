package models

import (
	"github.com/gofrs/uuid"
)

type Product struct {
	ID       *uuid.UUID `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"unique" json:"name"`
	Category uint       `json:"category"`

	Files []File `gorm:"many2many:product_files;"`
}

func NewProduct() *Product {
	ID, _ := uuid.NewV7()
	return &Product{ID: &ID}
}

type ProductFile struct {
	ProductID uint `json:"product_id"`
	FileID    uint `json:"file_id"`
}
