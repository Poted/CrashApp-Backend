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
