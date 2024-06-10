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

	// fmt.Printf("payServ.Countries: %v\n", payServ.Countries)
	fmt.Printf("payServ.locals: %v\n", payServ.Locals)
	fmt.Printf("payServ.locals.name: %v\n", payServ.Locals.Name)
	// fmt.Printf("payServ.Countries: %v\n", payServ.Currencies)

	// x, _ := payServ.Get().ByCountry(context.Background(), 2)

	// fmt.Printf("x: %v\n", x)

	// payServ.Set().ByCountry(context.Background(), "Brazylia", 2)

	// x, _ = payServ.Get().ByCountry(context.Background(), 2)

	// fmt.Printf("x: %v\n", x)

	handler.HttpClient()

}

func Hello() {

	serv := Server{
		TransformName: addPrefix,
	}

	serv.HandkeRequest("Whatever")

}

type Server struct {
	TransformName TransformFunc
}

type TransformFunc func(string) string

func addPrefix(input string) string {

	return input + " xD"
}

func (s *Server) HandkeRequest(fileName string) error {

	transName := s.TransformName(fileName)

	fmt.Printf("transName: %v\n", transName)

	return nil
}
