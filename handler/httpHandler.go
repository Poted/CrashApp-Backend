package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func HttpClient() {

	app := fiber.New()

	app.Use(cors.New())

	storageRouter(app)

	app.Listen(":80")

}

func storageRouter(app *fiber.App) {

	app.Get("/filesList", FilesList)

	app.Get("/getFileData/:id", GetFileData)

	app.Patch("/updateFile/:id", UpdateFileData)

	app.Get("/getFile/:id", GetFile)

	app.Post("/saveFile", SaveFile)

	app.Delete("/deleteFile/:id", DeleteFile)
}
