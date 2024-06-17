package internal

import (
	"context"
	"fmt"
	"go_app/backend/db"
	"go_app/backend/models"
	"sync"
)

type PaymentMethodService[T ILocals] struct {
	Locals *T
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

type ILocals interface {
	models.Country | models.Currency | models.ShippingMethod | models.Group
}

type GetByPreference func(context.Context, int, chan<- *models.Locale, *sync.WaitGroup) error

type SetUsingPreference func(context.Context, string, ...uint) error

func (pms *PaymentMethodService[T]) Get() *PaymentMethodServiceGetter {

	some := *pms.Locals
	_ = some

	someMore := &PaymentMethodServiceGetter{
		ByCountry:        getPaymentPreferenceByCountryService,
		ByCurrency:       getPaymentPreferenceByCurrencyService,
		ByShippingMethod: getPaymentPreferenceByShippingService,
		ByGroup:          getPaymentPreferenceByGroupService,
	}

	xd := someMore.ByCountry

	fmt.Printf("xd: %v\n", xd)

	// someMore.val = some

	return someMore
	// return nil
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

func (s *PaymentMethodService[T]) AddPrefix() {

}

func getPaymentPreferenceByCountryService(ctx context.Context, unitID int, ch chan<- *models.Locale, wg *sync.WaitGroup) error {

	defer wg.Done()

	country := models.Locale{}.Country

	err := db.Database.
		Where("id = ?", unitID).
		First(&country).
		Error

	select {
	// case ch <- country:
	case ch <- &models.Locale{Country: country}:
	case <-ctx.Done():
		fmt.Println("Fetch cancelled for ID:", unitID)
	}

	return err
}

func getPaymentPreferenceByCurrencyService(ctx context.Context, unitID int, ch chan<- *models.Locale, wg *sync.WaitGroup) error {

	currency := models.Currency{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&currency).
		Error

	select {
	case ch <- &models.Locale{Currency: currency}:
	case <-ctx.Done():
		fmt.Println("Fetch cancelled for ID:", unitID)
	}

	return err
}

func getPaymentPreferenceByShippingService(ctx context.Context, unitID int, ch chan<- *models.Locale, wg *sync.WaitGroup) error {

	method := models.ShippingMethod{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&method).
		Error

	select {
	case ch <- &models.Locale{ShippingMethod: method}:
	case <-ctx.Done():
		fmt.Println("Fetch cancelled for ID:", unitID)
	}

	return err
}

func getPaymentPreferenceByGroupService(ctx context.Context, unitID int, ch chan<- *models.Locale, wg *sync.WaitGroup) error {

	group := models.Group{}

	err := db.Database.
		Where("id = ?", unitID).
		First(&group).
		Error

	select {
	case ch <- &models.Locale{Group: group}:
	case <-ctx.Done():
		fmt.Println("Fetch cancelled for ID:", unitID)
	}

	return err
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

		e = db.Database.
			Create(&models.Country{Name: name}).
			Error

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
