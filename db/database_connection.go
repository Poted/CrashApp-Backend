package db

import (
	"errors"
	"go_app/backend/errorz"
	"go_app/backend/models"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	Database *gorm.DB
	dbOnce   sync.Once
)

func MigrateModels() {

	if Database != nil {
		err := Database.AutoMigrate(models.Storage()...)
		if err != nil {
			errorz.SendError(errors.New("cannot migrate storage"))
		}

		err = Database.AutoMigrate(models.Products()...)
		if err != nil {
			errorz.SendError(errors.New("cannot migrate products"))
		}

	}
}

func DBConnection() error {

	var err error

	dbOnce.Do(func() {

		connectionString := "sqlserver://sa:M1croshitSqlServer@localhost:1433?database=crashapp"

		Database, err = gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
		if err != nil {
			err = errorz.SendError(err)
		}

		Database = Database.Set("gorm:insert_option", "ON CONFLICT DO NOTHING")

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
