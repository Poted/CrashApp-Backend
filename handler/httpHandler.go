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

// ### Router ###

func router(app *fiber.App) {

	app.Get("/getIP", getIP)

	storageRouter(app)
	productsRouter(app)
}

// ### Storage Router ###

func storageRouter(app *fiber.App) {

	var storageHandler StorageHandler

	// Files:

	app.Post("/saveFile/:directory_id", storageHandler.SaveFile)
	app.Get("/filesList/:directory_id", storageHandler.FilesList)
	app.Get("/fileData/:file_id", storageHandler.GetFileData)
	app.Get("/file/:file_id", storageHandler.GetFile)

	// Directories:

	app.Post("/createFolder", storageHandler.CreateFolder)
	app.Get("/foldersList/:parentID", storageHandler.FoldersList)

}

// ### Product Router ###

func productsRouter(app *fiber.App) {

	app.Post("/products/create", CreateProduct)
	app.Get("/products/get/:id", GetProduct)
	app.Patch("/products/update/:id", UpdateProduct)
	app.Delete("/products/delete/:id", DeleteProduct)

	app.Get("/products/download/product_id/:id", DownloadProducts)

}

//																																				//

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

	return c.Status(ret.Code).Send(js)
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

	return c.Status(ferr.Code).Send(js)
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

	return c.Status(ret.Code).Send(js)
}
