package main

import (
	"go_app/backend/db"
	"go_app/backend/errorz"
	"go_app/backend/handler"
)

func main() {

	err := db.DBConnection()
	if err != nil {
		errorz.SendError(err)
	}
	defer db.CloseDB(db.Database)

	db.MigrateModels()

	handler.HttpClient()

}
