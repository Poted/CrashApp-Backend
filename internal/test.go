package internal

import (
	"context"
	"go_app/backend/db"
	"go_app/backend/models"
)

type PaymentMethodService struct {
	Countries       []models.Country
	Currencies      []models.Currency
	ShippingMethods []models.ShippingMethod
	Groups          []models.Group
}

type PaymentMethodServiceGetter struct {
	ByCountry        GetByPreference
	ByCurrency       GetByPreference
	ByShippingMethod GetByPreference
	ByGroup          GetByPreference
}

type PaymentMethodServiceSetter struct {
	ByCountry        SetUsingPreference
	ByCurrency       SetUsingPreference
	ByShippingMethod SetUsingPreference
	ByGroup          SetUsingPreference
}

type GetByPreference func(context.Context, int) (string, error)

type SetUsingPreference func(context.Context, string, ...uint) error

func (*PaymentMethodService) Get() *PaymentMethodServiceGetter {
	return &PaymentMethodServiceGetter{
		ByCountry:        getPaymentPreferenceByCountryService,
		ByCurrency:       getPaymentPreferenceByCurrencyService,
		ByShippingMethod: getPaymentPreferenceByShippingService,
		ByGroup:          getPaymentPreferenceByGroupService,
	}
}

func (*PaymentMethodService) Set() *PaymentMethodServiceSetter {
	return &PaymentMethodServiceSetter{
		ByCountry:        setPaymentPreferenceByCountryService,
		ByCurrency:       setPaymentPreferenceByCurrencyService,
		ByShippingMethod: setPaymentPreferenceByShippingService,
		ByGroup:          setPaymentPreferenceByGroupService,
	}
}

//	### GETTERS ###

func getPaymentPreferenceByCountryService(ctx context.Context, unitID int) error {

	// ### TO DO: ###
	// func getPaymentPreferenceByCountryService(ctx context.Context, unitID ...int) (string, error) {
	// Make a func that will take zero parameters: returning all items
	// one parameter: one object
	// multiple IDS: objects with given IDs

	// Make those functions just for retrieving data from arrays/slices and make a db call when it's needed

	country := models.Country{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&country).
		Error

	return err
}

func getPaymentPreferenceByCurrencyService(ctx context.Context, unitID int) (string, error) {

	currency := models.Country{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&currency).
		Error

	return currency.Name, err
}

func getPaymentPreferenceByShippingService(ctx context.Context, unitID int) (string, error) {

	method := models.Country{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&method).
		Error

	return method.Name, err
}

func getPaymentPreferenceByGroupService(ctx context.Context, unitID int) (string, error) {

	group := models.Country{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&group).
		Error

	return group.Name, err
}

//	### SETTERS ###

func setPaymentPreferenceByCountryService(ctx context.Context, name string, paymentMethodIDs ...uint) (e error) {

	if len(paymentMethodIDs) == 1 {

		e = db.Database.
			Model(&models.Country{}).
			Where("id = ?", paymentMethodIDs[0]).
			Updates(&models.Country{Name: name}).
			Error

	} else if len(paymentMethodIDs) > 1 {

		e = db.Database.
			Model(&models.Country{}).
			Where("id IN (?)", paymentMethodIDs).
			Updates(&models.Country{Name: name}).
			Error

	} else {
		// countries = append(countries, name)
	}

	return e
}

func setPaymentPreferenceByCurrencyService(ctx context.Context, name string, paymentMethodIDs ...uint) (e error) {

	if len(paymentMethodIDs) == 1 {

		e = db.Database.
			Model(&models.Currency{}).
			Where("id = ?", paymentMethodIDs[0]).
			Updates(&models.Currency{Name: name}).
			Error

	} else if len(paymentMethodIDs) > 1 {

		e = db.Database.
			Model(&models.Currency{}).
			Where("id IN (?)", paymentMethodIDs).
			Updates(&models.Currency{Name: name}).
			Error

	} else {
		// currencies = append(currencies, name)
	}

	return e
}

func setPaymentPreferenceByShippingService(ctx context.Context, name string, paymentMethodIDs ...uint) (e error) {

	if len(paymentMethodIDs) == 1 {

		e = db.Database.
			Model(&models.ShippingMethod{}).
			Where("id = ?", paymentMethodIDs[0]).
			Updates(&models.ShippingMethod{Name: name}).
			Error

	} else if len(paymentMethodIDs) > 1 {

		e = db.Database.
			Model(&models.ShippingMethod{}).
			Where("id IN (?)", paymentMethodIDs).
			Updates(&models.ShippingMethod{Name: name}).
			Error

	} else {
		// methods = append(methods, name)
	}

	return e
}

func setPaymentPreferenceByGroupService(ctx context.Context, name string, paymentMethodIDs ...uint) (e error) {

	if len(paymentMethodIDs) == 1 {

		e = db.Database.
			Model(&models.Group{}).
			Where("id = ?", paymentMethodIDs[0]).
			Updates(&models.Group{Name: name}).
			Error

	} else if len(paymentMethodIDs) > 1 {

		e = db.Database.
			Model(&models.Group{}).
			Where("id IN (?)", paymentMethodIDs).
			Updates(&models.Group{Name: name}).
			Error

	} else {
		// countries = append(countries, name)
	}

	return e
}
