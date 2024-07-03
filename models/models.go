package models

func List(items ...interface{}) []interface{} {
	return items
}

func StorageList() []interface{} {
	return List(
		Directory{},
		File{},
	)
}

func ProductsList() []interface{} {
	return List(
		Product{},
	)
}

func LocaleList() []interface{} {
	return List(
		Country{},
		Currency{},
		ShippingMethod{},
		Group{},
		ShippingCountry{},
	)
}

type Locale struct {
	Country
	Currency
	ShippingMethod
	Group
}
