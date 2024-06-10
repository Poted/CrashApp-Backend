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

	payServ := internal.PaymentMethodService[models.Country]{}

	err = payServ.Get().ByCountry(context.Background(), 2)
	if err != nil {
		errorz.SendError(err)
	}

	err = payServ.Get().ByCurrency(context.Background(), 2)
	if err != nil {
		errorz.SendError(err)
	}

	// payServ.Locals = models.Country{
	// 	Name: "POLSKA GUROM ",
	// }

	fmt.Printf("payServ.Countries: %v\n", payServ.Countries)
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
