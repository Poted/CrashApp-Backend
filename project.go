package main

import (
	"context"
	"fmt"
	"go_app/backend/db"
	"go_app/backend/errorz"
	"go_app/backend/handler"
	"go_app/backend/internal"
	"go_app/backend/models"
)

func main() {

	err := db.DBConnection()
	if err != nil {
		errorz.SendError(err)
	}
	defer db.CloseDB(db.Database)

	db.MigrateModels()

	payServ := internal.PaymentMethodService[models.Country]{
		Locals: &models.Country{
			Name: "URUGWAJ",
		},
	}

	err = payServ.Get().ByCountry(context.Background(), 2)
	if err != nil {
		errorz.SendError(err)
	}

	err = payServ.Get().ByCurrency(context.Background(), 2)
	if err != nil {
		errorz.SendError(err)
	}

	fmt.Printf("payServ.currency: %v\n", payServ.Locals.Name)
	payServ.Locals = &models.Country{
		Name: "POLSKA GUROM",
	}

	fmt.Printf("payServ.locals: %v\n", payServ.Locals)
	fmt.Printf("payServ.locals.name: %v\n", payServ.Locals.Name)

	payServ.Set().ByCountry(context.Background(), "Tajlandia")
	err = payServ.Set().ByCountry(context.Background(), "Mozambik", 2)
	// err = payServ.Set().ByCountry(context.Background(), "Mozambik", 1, 2)
	if err != nil {
		errorz.SendError(err)
	}

	err = payServ.Get().ByCountry(context.Background(), 2)
	if err != nil {
		errorz.SendError(err)
	}

	fmt.Printf("payServ.locals.name: %v\n", payServ.Locals.Name)

	fmt.Printf("x: %v\n", err)

	handler.HttpClient()

}
