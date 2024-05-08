package main

import (
	"errors"
	"go_app/backend/errorz"
	"go_app/backend/handler"
	"go_app/backend/models"
	"sync"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	db     *gorm.DB
	dbOnce sync.Once
)

func main() {

	err := DBConnection()
	if err != nil {
		errorz.SendError(err)
	}
	defer CloseDB(db)

	handler.HttpClient()

}

func MigrateModels() error {
	if db != nil {
		err := db.AutoMigrate(models.List()...)
		if err != nil {
			return errorz.SendError(err)
		}
	}
	return nil
}

func DBConnection() error {

	// docker run -e 'ACCEPT_EULA=Y' -e 'SA_PASSWORD=M1croshitSqlServer' -p 1433:1433 --name sql_server_container -d mcr.microsoft.com/mssql/server:latest

	var err error

	dbOnce.Do(func() {

		connectionString := "sqlserver://sa:M1croshitSqlServer@localhost:1433?database=crashapp"

		db, err = gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
		if err != nil {
			err = errorz.SendError(err)
		}

		db = db.Set("gorm:insert_option", "ON CONFLICT DO NOTHING")

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
