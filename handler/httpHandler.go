package handler

import (
	"encoding/json"
	"fmt"
	"go_app/backend/errorz"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	fiber_recover "github.com/gofiber/fiber/v2/middleware/recover"
)

func HttpClient() {

	app := fiber.New(fiber.Config{
		ProxyHeader: fiber.HeaderXForwardedFor,
	})
	app.Use(fiber_recover.New(
		fiber_recover.Config{
			EnableStackTrace: true,
			StackTraceHandler: func(c *fiber.Ctx, e interface{}) {
				fmt.Printf("debug.Stack(): %v\n", string(debug.Stack()))
			},
		},
	))
	app.Use(cors.New())
	router(app)

	app.Listen(":80")

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
	app.Get("/products/get/:id", GetProduct)
	app.Patch("/products/update/:id", UpdateProduct)
	app.Delete("/products/delete/:id", DeleteProduct)

}

type Response struct {
	Type    ResponseType `json:"type"`
	Message string       `gorm:"required" json:"message"`
	Data    *interface{} `json:"data,omitempty"`
}

type Status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func NewStatus(code int, message ...string) *Status {

	return &Status{
		Code: code,
		Message: func() string {
			if len(message) > 0 {
				return message[0]
			}
			return ""
		}(),
	}
}

type ResponseType string

const (
	StatusSuccess ResponseType = "success"
	StatusError   ResponseType = "error"
	StatusInfo    ResponseType = "info"
)

func SuccessResponse(c *fiber.Ctx, ret *Status, data interface{}) error {

	js, err := json.Marshal(Response{
		Type:    StatusSuccess,
		Message: ret.Message,
		Data:    &data,
	})
	if err != nil {
		return errorz.SendError(err)
	}

	c.Status(ret.Code).Send(js)

	return nil
}

func ErrorResponse(c *fiber.Ctx, ferr *fiber.Error, data *interface{}) error {

	fmt.Print(errorz.SendError(ferr))

	js, err := json.Marshal(Response{
		Type:    StatusError,
		Message: ferr.Message,
		Data:    data,
	})
	if err != nil {
		return errorz.SendError(err)
	}

	c.Status(ferr.Code).Send(js)

	return nil
}

func CreateInfoResponse(c *fiber.Ctx, ret *Status, data interface{}) error {

	js, err := json.Marshal(Response{
		Type:    StatusInfo,
		Message: ret.Message,
		Data:    &data,
	})
	if err != nil {
		return errorz.SendError(err)
	}

	c.Status(ret.Code).Send(js)

	return nil
}
