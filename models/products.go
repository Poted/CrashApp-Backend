package models

import (
	"github.com/gofrs/uuid"
)

type Product struct {
	ID       *uuid.UUID `gorm:"primaryKey" json:"id"`
	Name     string     `gorm:"unique" json:"name"`
	Category uint       `json:"category"`
}

func NewProduct() *Product {

	ID, _ := uuid.NewV7()

	if ID != uuid.Nil {
		return &Product{ID: &ID}
	}
	return NewProduct()
}
