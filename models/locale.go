package models

import "gorm.io/gorm"

type Country struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Currency struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type ShippingMethod struct {
	gorm.Model
	Name string `gorm:"unique"`
}

type Group struct {
	gorm.Model
	Name string `gorm:"unique"`
}
