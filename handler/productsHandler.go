package handler

import (
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
		return ErrorResponse(c, fiber.NewError(400, "wrong input; body is invalid"), nil)
	}

	createdProduct, err := productRepo.CreateProduct(ctx, model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, "wrong input; query params are invalid"), nil)
	}

	return SuccessResponse(c, NewStatus(200, "successfully created a product"), createdProduct)
}

func GetProduct(c *fiber.Ctx) error {

	ctx := c.UserContext()
	id := uuid.FromStringOrNil(c.Params("id"))

	if id == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "wrong input; ID is invalid"), nil)
	}
	model := models.Product{ID: &id}

	err := c.BodyParser(&model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	product, err := productRepo.GetProduct(ctx, &model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(500, err.Error()), nil)
	}

	return SuccessResponse(c, NewStatus(200, "successfully fetched product data"), product)
}

func UpdateProduct(c *fiber.Ctx) error {

	ctx := c.UserContext()
	id := uuid.FromStringOrNil(c.Params("id"))

	if id == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "wrong input; ID missing"), nil)
	}
	model := models.Product{ID: &id}

	err := c.BodyParser(&model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	updatedProduct, err := productRepo.UpdateProduct(ctx, &model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(500, err.Error()), nil)
	}

	return SuccessResponse(c, NewStatus(200, "successfully updated a product"), updatedProduct)
}

func DeleteProduct(c *fiber.Ctx) error {

	ctx := c.UserContext()
	id := uuid.FromStringOrNil(c.Params("id"))

	if id == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "wrong input; ID missing"), nil)
	}
	model := models.Product{ID: &id}

	err := c.BodyParser(&model)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	ferr := productRepo.DeleteProduct(ctx, &model)
	if ferr != nil {
		return ErrorResponse(c, ferr, nil)
	}

	return SuccessResponse(c, NewStatus(201, "successfully deleted a product"), nil)
}
