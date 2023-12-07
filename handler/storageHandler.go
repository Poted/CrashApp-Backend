package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_app/backend/errorz"
	"go_app/backend/helperz"
	"go_app/backend/internal"
	"go_app/backend/models"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var il = internal.New()

func SaveFile(c *fiber.Ctx) error {

	// Get file from request body
	fileMultipart, err := c.FormFile("file_name")
	if err != nil {
		errorz.SendError(err)
	}

	fileModel := models.File{
		ID:        uuid.New(),
		Name:      fileMultipart.Filename,
		Size:      fileMultipart.Size,
		Directory: "",
	}

	jsonDB, err := os.OpenFile("C:/Users/ojpkm/Documents/go_app/Database/files.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonDB.Close()

	// Decode the existing JSON data from the file
	var files []models.File
	decoder := json.NewDecoder(jsonDB)
	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	// Append the new record to the existing data
	files = append(files, fileModel)

	// Seek to the beginning of the file to overwrite the existing data
	if _, err := jsonDB.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Create a JSON encoder and encode the updated data to the file
	encoder := json.NewEncoder(jsonDB)
	if err := encoder.Encode(files); err != nil {
		log.Fatal(err)
	}

	// Save file to storage folder
	err = c.SaveFile(fileMultipart, fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/storage/%s", fileModel.ID))
	if err != nil {
		return err
	}

	return c.Status(201).SendString("Created")
}

func FilesList(c *fiber.Ctx) error {

	storage := internal.New()

	files, err := storage.ReadFilesList()
	if err != nil {
		return err
	}

	formatedData, err := json.Marshal(files)
	if err != nil {
		return err
	}

	return c.Send(formatedData)
}

func GetFileData(c *fiber.Ctx) error {

	id := c.Params("id")
	storage := internal.New()
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	fileData, err := storage.ReadFileData(&uuid)
	if err != nil {
		return err
	}

	file, err := json.Marshal(fileData)
	if err != nil {
		return err
	}

	return c.Send(file)

}

func UpdateFileData(c *fiber.Ctx) error {

	id := c.Params("id")
	storage := internal.New()
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	var fileModel models.File
	err = c.BodyParser(&fileModel)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "wrong data provided")
	}

	fileData, err := storage.UpdateFileData(&uuid, &fileModel)
	if err != nil {
		return err
	}

	return c.Status(200).Send(*fileData)

}

func GetFile(c *fiber.Ctx) error {

	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	storage := il
	fileData, err := storage.ReadFileData(&uuid)
	if err != nil {
		return err
	}

	file, err := os.ReadFile(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Storage/%s", fileData.ID))
	if err != nil {
		return err
	}

	if file == nil {
		return errors.New("cannot send a file")
	}

	if strings.Contains(fileData.Name, ".pdf") {
		c.Set("Content-Type", "application/pdf")
	} else {
		c.Set("Content-Type", "image/jpeg")
	}

	c.Set("Content-Length", fmt.Sprint(file))

	return c.Status(fiber.StatusOK).Send(file)
}

func DeleteFile(c *fiber.Ctx) error {

	id := c.Params("id")
	uuid, err := uuid.Parse(id)
	if err != nil {
		return c.Status(503).Send([]byte("cannot find a file"))
	}

	storage := il
	err = storage.DeleteFile(&uuid)
	if err != nil {
		return c.Status(503).Send([]byte("cannot find a file"))
	}

	return c.Status(204).Send([]byte("Succesfully removed file"))

}

func CreateFolder(c *fiber.Ctx) error {

	body := models.Directory{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err = helperz.DataBaseInsert(body, "directory")
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	err = helperz.DataBaseUpdate(body, "directory", &body.ID)
	if err != nil {
		return c.Status(500).SendString(err.Error())
	}

	return c.Status(200).SendString("Created")
}