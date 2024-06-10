package internal

import (
	"context"
	"fmt"
	"go_app/backend/db"
	"go_app/backend/models"
)

type PaymentMethodService[T ILocals] struct {
	Countries       []models.Country
	Currencies      []models.Currency
	ShippingMethods []models.ShippingMethod
	Groups          []models.Group
	Locals          T
}

type PaymentMethodServiceGetter struct {
	ByCountry        GetByPreference
	ByCurrency       GetByPreference
	ByShippingMethod GetByPreference
	ByGroup          GetByPreference
}

type PaymentMethodServiceSetter struct {

	// ByCountry(context.Context, name, IDs) error
	ByCountry SetUsingPreference
	// ByCurrency(context.Context, name, IDs) error
	ByCurrency SetUsingPreference
	// ByShipping(context.Context, name, IDs) error
	ByShippingMethod SetUsingPreference
	// ByGroup(context.Context, name, IDs) error
	ByGroup SetUsingPreference
}

type ILocals interface {
	models.Country | models.Currency | models.ShippingMethod | models.Group
}

type GetByPreference func(context.Context, int) error

// type GetByPreference func(context.Context, int) error
// type GetByPreference func(context.Context, int) (string, error)

type SetUsingPreference func(context.Context, string, ...uint) error

func (pms *PaymentMethodService[T]) Get() *PaymentMethodServiceGetter {

	pmsg := &PaymentMethodServiceGetter{
		ByCountry:  getPaymentPreferenceByCountryService,
		ByCurrency: getPaymentPreferenceByCurrencyService,
		// ByShippingMethod: getPaymentPreferenceByShippingService,
		// ByGroup:          getPaymentPreferenceByGroupService,

	}

	some := PaymentMethodService[T]{
		Countries: []models.Country{
			{
				Name: "Greenland",
			},
		},
	}

	fmt.Printf("some: %v\n", some)

	pms.Countries = []models.Country{
		{
			Name: "MERICA FUCK YEAH",
		},
	}

	return pmsg
}

func (*PaymentMethodService[T]) Set() *PaymentMethodServiceSetter {
	return &PaymentMethodServiceSetter{
		ByCountry:        setPaymentPreferenceByCountryService,
		ByCurrency:       setPaymentPreferenceByCurrencyService,
		ByShippingMethod: setPaymentPreferenceByShippingService,
		ByGroup:          setPaymentPreferenceByGroupService,
	}
}

//	### GETTERS ###

func getPaymentPreferenceByCountryService(ctx context.Context, unitID int) error {

	// func getPaymentPreferenceByCountryService(ctx context.Context, unitID int) error {

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
	// return country.Name, err
}

func getPaymentPreferenceByCurrencyService(ctx context.Context, unitID int) error {

	fmt.Printf("\"Hiii\": %v\n", "second function")

	currency := models.Country{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&currency).
		Error

	return err
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
