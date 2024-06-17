package db

import (
	"errors"
	"go_app/backend/errorz"
	"go_app/backend/models"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	Database *gorm.DB
	dbOnce   sync.Once
)

func MigrateModels() {

	if Database != nil {
		err := Database.AutoMigrate(models.StorageList()...)
		if err != nil {
			errorz.SendError(errors.New("cannot migrate storage"))
		}

		err = Database.AutoMigrate(models.ProductsList()...)
		if err != nil {
			errorz.SendError(errors.New("cannot migrate products"))
		}

		err = Database.AutoMigrate(models.LocaleList()...)
		if err != nil {
			errorz.SendError(errors.New("cannot migrate products"))
		}

		// populateCountries(Database)
		// populateCurrencies(Database)
		// populateShippingMethods(Database)
		// populateGroups(Database)

	}
}

func DBConnection() error {

	var err error

	dbOnce.Do(func() {

		connectionString := "sqlserver://sa:M1croshitSqlServer@localhost:1433?database=crashapp"

		Database, err = gorm.Open(sqlserver.Open(connectionString), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Error),
		})
		if err != nil {
			err = errorz.SendError(err)
		}

	})

	return err
}

func CloseDB(db *gorm.DB) error {

	if db != nil {
		DBobject, err := db.DB()
		if err == nil {
			err = DBobject.Close()
		}
		return err
	}
	return errors.New("no db found")
}

// Function to populate the countries table
func populateCountries(db *gorm.DB) {
	countries := []string{"Szwecja", "Estornia", "Bułgaria"}

	for _, name := range countries {
		country := models.Country{Name: name}
		db.Create(&country)
	}
}

// Function to populate the currencies table
func populateCurrencies(db *gorm.DB) {
	currencies := []string{"Korona", "Zloty", "Euro"}

	for _, name := range currencies {
		currency := models.Currency{Name: name}
		db.Create(&currency)
	}
}

// Function to populate the shippingMethods table
func populateShippingMethods(db *gorm.DB) {
	shippingMethods := []string{"Chiny", "Inpost", "DHL"}

	for _, name := range shippingMethods {
		method := models.ShippingMethod{Name: name}
		db.Create(&method)
	}
}

// Function to populate the groups table
func populateGroups(db *gorm.DB) {
	groups := []string{"Klienci", "Goscie", "VIP"}

	for _, name := range groups {
		group := models.Group{Name: name}
		db.Create(&group)
	}
}
