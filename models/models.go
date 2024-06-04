package models

func List(category ...interface{}) []interface{} {
	return category
}

func Storage() []interface{} {
	return List(
		Directory{},
		File{},
	)
}

func Products() []interface{} {
	return List(
		Product{},
	)
}

func Locale() []interface{} {
	return List(
		Country{},
		Currency{},
		ShippingMethod{},
		Group{},
	)
}
