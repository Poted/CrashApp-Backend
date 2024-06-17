package main

import (
	"context"
	"fmt"
	"go_app/backend/db"
	"go_app/backend/errorz"
	"go_app/backend/handler"
	"go_app/backend/internal"
	"go_app/backend/models"
	"sync"
)

func main() {

	err := db.DBConnection()
	if err != nil {
		errorz.SendError(err)
	}
	defer db.CloseDB(db.Database)

	db.MigrateModels()

	payServ, hell := newFunction(err)

	fmt.Printf("hell: %v\n", hell)

	fmt.Printf("payServ.Locals.Name: %v\n", payServ.Locals.Name)
	payServ.Locals = &models.Country{
		Name: "Second main",
	}

	fmt.Printf("payServ.locals: %v\n", payServ.Locals)
	fmt.Printf("payServ.locals.name: %v\n", payServ.Locals.Name)

	// payServ.Set().ByCountry(context.Background(), "set Tajlandia")
	// err = payServ.Set().ByCountry(context.Background(), "set Georgia", 67)
	// if err != nil {
	// 	errorz.SendError(err)
	// }

	handler.HttpClient()

}

func newFunction(err error) (internal.PaymentMethodService[models.Country], *models.Locale) {

	countryChan := make(chan *models.Locale)
	var wg sync.WaitGroup

	payServ := internal.PaymentMethodService[models.Country]{
		Locals: &models.Country{
			Name: "first main",
		},
	}

	for i := 1; i <= 4; i++ {
		wg.Add(1)
		go func(ctx context.Context, i int, ch chan *models.Locale, wg *sync.WaitGroup) {

			err = payServ.Get().ByCountry(context.Background(), 67, ch, wg)
			if err != nil {
				errorz.SendError(err)
			}

		}(context.Background(), i, countryChan, &wg)
	}

	hell := <-countryChan
	return payServ, hell
}
