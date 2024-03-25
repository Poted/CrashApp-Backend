package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_app/backend/errorz"
	"go_app/backend/internal"
	"go_app/backend/models"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
)

// var il = internal.New()

var storage = internal.New()

type GeoIPResponse struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code"`
	RegionName  string  `json:"region_name"`
	City        string  `json:"city"`
	ZipCode     string  `json:"zip_code"`
	TimeZone    string  `json:"time_zone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
	MetroCode   int     `json:"metro_code"`
}

func getIP(c *fiber.Ctx) error {

	ip := c.IP() // Example IP address, you can replace it with the actual IP

	// Fetch GeoIP information from the API
	resp, err := http.Get("https://freegeoip.app/json/" + ip)
	if err != nil {
		fmt.Println("Error fetching GeoIP data:", err)
		// return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		// return
	}

	// Parse the JSON response
	var geoIP GeoIPResponse
	err = json.Unmarshal(body, &geoIP)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		// return
	}

	some := map[string]interface{}{
		"ip":           geoIP.IP,
		"country_code": geoIP.CountryCode,
		"country_name": geoIP.CountryName,
		"region_code":  geoIP.RegionCode,
		"region_name":  geoIP.RegionName,
		"city":         geoIP.City,
		"zip_code":     geoIP.ZipCode,
		"time_zone":    geoIP.TimeZone,
		"latitude":     geoIP.Latitude,
		"longitude":    geoIP.Longitude,
		"metro_code":   geoIP.MetroCode,
	}

	_ = some

	// return c.JSON(fiber.Map{
	// 	"ip":           geoIP.IP,
	// 	"country_code": geoIP.CountryCode,
	// 	"country_name": geoIP.CountryName,
	// 	"region_code":  geoIP.RegionCode,
	// 	"region_name":  geoIP.RegionName,
	// 	"city":         geoIP.City,
	// 	"zip_code":     geoIP.ZipCode,
	// 	"time_zone":    geoIP.TimeZone,
	// 	"latitude":     geoIP.Latitude,
	// 	"longitude":    geoIP.Longitude,
	// 	"metro_code":   geoIP.MetroCode,
	// })

	somefunc(c)

	return nil

}

// ////////////////////
func somefunc(c *fiber.Ctx) {

	accountHandler := &AccountHandler{AccountNotifier: SomeNotifier{}}

	accountHandler.handleCreateAccount(c)

}

type AccountNotifier interface {
	NotifyAccountCreated(context.Context, Account) error
}

type SomeNotifier struct{}
type AccountHandler struct {
	AccountNotifier AccountNotifier
}

func (s SomeNotifier) NotifyAccountCreated(ctx context.Context, account Account) error {
	fmt.Printf("Created %s, %s", account.Email, account.Username)
	return nil
}

type Account struct {
	Username string
	Email    string
}

func (h *AccountHandler) handleCreateAccount(c *fiber.Ctx) {

	h.AccountNotifier.NotifyAccountCreated(c.Context(), Account{Username: "Some", Email: "XDSS@DSS.CJJ"})

	c.Send(bytes.NewBufferString("HIXD").Bytes())
}

// /////////////////////
func SaveFile(c *fiber.Ctx) error {

	// Get file from request body
	fileMultipart, err := c.FormFile("file_name")
	if err != nil {
		errorz.SendError(err)
	}

	directory_id := c.Params("directory_id")

	generateID, err := uuid.NewV7()
	if err != nil {
		return errorz.SendError(err)
	}

	fileModel := models.File{
		ID:   generateID,
		Name: fileMultipart.Filename,
		Size: fileMultipart.Size,
		Directory: func() string {
			if directory_id == "1" {
				folder, err := storage.ReadFolderData(nil, "main")
				if err != nil {
					return ""
				}
				return folder.ID.String()
			}
			return directory_id
		}(),
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

	files, err := storage.ReadFilesList()
	if err != nil {
		return errorz.SendError(err)
	}

	id := c.Params("directory_id")

	if id == "1" {

		folder, err := storage.ReadFolderData(nil, "main")
		if err != nil {
			return errorz.SendError(err)
		}
		id = folder.ID.String()
	}

	output := []models.File{}

	for _, v := range files {
		if v.Directory == id {
			output = append(output, v)
		}
	}

	formatedData, err := json.Marshal(output)
	if err != nil {
		return errorz.SendError(err)
	}

	return c.Send(formatedData)
}

func GetFileData(c *fiber.Ctx) error {

	id := c.Params("id")
	storage := internal.New()
	uuid, err := uuid.FromString(id)
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
	uuid, err := uuid.FromString(id)
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
	uuid, err := uuid.FromString(id)
	if err != nil {
		return err
	}

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
	uuid, err := uuid.FromString(id)
	if err != nil {
		return c.Status(503).Send([]byte("cannot find a file"))
	}

	err = storage.DeleteFile(&uuid)
	if err != nil {
		return c.Status(503).Send([]byte("cannot find a file"))
	}

	return c.Status(204).Send([]byte("Succesfully removed file"))

}

func DeleteFolder(c *fiber.Ctx) error {

	id := c.Params("id")
	uuid, err := uuid.FromString(id)
	if err != nil {
		return c.Status(503).SendString(err.Error())
	}

	err = storage.DeleteFolder(&uuid)
	if err != nil {
		return c.Status(503).Send([]byte("cannot find a folder"))
	}

	return c.Status(204).Send([]byte("Succesfully removed a folder"))

}

func SaveFolder(c *fiber.Ctx) error {

	body := models.Directory{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	_, err = storage.CreateFolder(&body)
	if err != nil {
		return errorz.SendError(err)
	}

	return c.Status(201).Send([]byte("Succesfully created a file"))
}

func GetFolders(c *fiber.Ctx) error {

	folders, err := storage.GetFolders()
	if err != nil {
		return err
	}

	return c.Status(200).JSON(folders)
}

func EditFolder(c *fiber.Ctx) error {

	body := models.Directory{}

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	err = storage.UpdateFolderData(&body.ID, &body)
	if err == nil {
		return c.Status(204).Send([]byte("Succesfully edited a file"))
	}

	return nil
}
