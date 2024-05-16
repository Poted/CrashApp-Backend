package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"go_app/backend/errorz"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func HttpClient() {

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("debug.Stack(): %v\n", debug.Stack())
		}
	}()

	app := fiber.New()
	app.Use(recoverMiddleware)
	app.Use(cors.New())
	router(app)
	app.Listen(":80")

}

func recoverMiddleware(c *fiber.Ctx) error {
	defer func() {
		if r := recover(); r != nil {
			// Recover from panic and respond with an error message
			c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal Server Error",
			})
		}
	}()
	// Continue to the next middleware or handler
	return c.Next()
}

func router(app *fiber.App) {
	storageRouter(app)
	productsRouter(app)
}

func storageRouter(app *fiber.App) {

	app.Get("/getIP", getIP)

	// Files:

	app.Get("/filesList/:directory_id", FilesList)
	app.Get("/getFileData/:id", GetFileData)
	app.Patch("/updateFile/:id", UpdateFileData)
	app.Get("/getFile/:id", GetFile)
	app.Post("/saveFile/:directory_id", SaveFile)
	app.Delete("/deleteFile/:id", DeleteFile)

	// Directories:

	app.Post("/SaveFolder", SaveFolder)
	app.Get("/GetFolders", GetFolders)
	app.Put("/EditFolder", EditFolder)
	app.Delete("/DeleteFolder/:id", DeleteFolder)

}

func productsRouter(app *fiber.App) {

	app.Post("/products/create", CreateProduct)
	app.Patch("/products/update/:id", UpdateProduct)

}

type Response struct {
	Type    ResponseType `json:"type"`
	Message string       `gorm:"required" json:"message"`
	Data    *interface{} `json:"data,omitempty"`
}

type ResponseType string

const (
	StatusSuccess ResponseType = "success"
	StatusError   ResponseType = "error"
	StatusInfo    ResponseType = "info"
)

func CreateErrorResponse(c *fiber.Ctx, ferr *fiber.Error, data *interface{}) {

	fmt.Printf("errorz.SendError(ferr): %v\n", errorz.SendError(ferr))

	js, err := json.Marshal(Response{
		Type:    StatusError,
		Message: ferr.Message,
		Data:    data,
	})
	if err != nil {
		errorz.SendError(err)
	}

	c.Status(ferr.Code).Send(js)

}

func CreateInfoResponse(ctx context.Context, message string, data interface{}) Response {
	return Response{
		Type:    StatusInfo,
		Message: message,
		Data:    &data,
	}
}

func CreateSuccessResponse(ctx context.Context, message string, data interface{}) Response {
	return Response{
		Type:    StatusSuccess,
		Message: message,
		Data:    &data,
	}
}
