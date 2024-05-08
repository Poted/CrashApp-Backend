package models

import "fmt"

func List(category ...interface{}) []interface{} {
	return category
}

type ModelCar []interface{}

func Storage() []interface{} {

	var storage ModelCar = []interface{}{
		Directory{},
		File{},
	}

	return List(storage...)
}

func test() {

	fmt.Printf("Storage(): %v\n", Storage())

}
