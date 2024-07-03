package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	Name string `gorm:"unique" json:"name"`

	ShippingMethods []ShippingMethod `gorm:"many2many:shipping_countries;"`
}

type Currency struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type ShippingMethod struct {
	gorm.Model
	Name string `gorm:"unique"`

	Countries []Country `gorm:"many2many:shipping_countries;"`
}

type Group struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type ShippingCountry struct {
	CountryID        uint `gorm:"primaryKey"`
	ShippingMethodID uint `gorm:"primaryKey"`
}
