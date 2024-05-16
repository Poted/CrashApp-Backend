package handler

import (
	"fmt"
	"go_app/backend/internal"
	"go_app/backend/models"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

var productRepo = internal.NewProducts()

func CreateProduct(c *fiber.Ctx) error {

	ctx := c.UserContext()

	model := models.NewProduct()

	err := c.BodyParser(&model)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	createdProduct, err := productRepo.CreateProduct(ctx, model)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	fmt.Printf("createdProduct: %v\n", createdProduct)

	return nil
}

func UpdateProduct(c *fiber.Ctx) error {

	ctx := c.UserContext()
	id := uuid.FromStringOrNil(c.Params("id"))

	if id == uuid.Nil {

		CreateErrorResponse(c, fiber.NewError(400, "wrong input; ID missing"), nil)
	}
	model := models.Product{}

	err := c.BodyParser(&model)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	updatedProduct, err := productRepo.UpdateProduct(ctx, &model)
	if err != nil {
		return fiber.NewError(400, err.Error())
	}

	fmt.Printf("updatedProduct: %v\n", updatedProduct)

	return nil
}
