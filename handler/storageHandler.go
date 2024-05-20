package handler

import (
	"errors"
	"fmt"
	"go_app/backend/internal"
	"go_app/backend/models"
	"mime/multipart"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofrs/uuid"
	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

type IStorageHandler interface {
	SaveFile(c *fiber.Ctx) error
	FilesList(c *fiber.Ctx) error
	GetFileData(c *fiber.Ctx) error
	GetFile(c *fiber.Ctx) error
	CreateFolder(c *fiber.Ctx) error
	FoldersList(c *fiber.Ctx) error
}

type StorageHandler struct {
	internal.IStorageInternal
}

func NewStorageHandler() IStorageHandler {
	return &StorageHandler{
		internal.NewStorage(),
	}
}

// var h.IStorageInternal = internal.NewStorage()

func (h *StorageHandler) SaveFile(c *fiber.Ctx) error {

	fileMultipart, err := c.FormFile("file_name")
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	extention, err := readFileExtention(fileMultipart)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if !strings.Contains(extention.MIME.Value, "image/") {
		return ErrorResponse(c, fiber.NewError(400, fmt.Sprintf("cannot save file of type %s", extention.Extension)), nil)
	}

	if fileMultipart.Size > 250000 {
		return ErrorResponse(c, fiber.NewError(400, "too big file"), nil)
	}

	directory_id := c.Params("directory_id")

	fileModel := models.File{
		Name:      fileMultipart.Filename,
		Size:      fileMultipart.Size,
		Directory: directory_id,
	}.NewFile()

	saveFileFunc := func(path string) error {
		return c.SaveFile(fileMultipart, path)
	}

	fileModel, err = h.IStorageInternal.SaveFile(fileModel, saveFileFunc)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	return SuccessResponse(c, NewStatus(200, "Created"), fileModel)
}

func readFileExtention(fileMultipart *multipart.FileHeader) (*types.Type, error) {

	file, err := fileMultipart.Open()
	if err != nil {
		return nil, err
	}
	defer func() {
		file.Seek(0, 0)
		file.Close()
	}()

	// Read the first 261 bytes (enough for most file type detections)
	head := make([]byte, 261)
	_, err = file.Read(head)
	if err != nil {
		return nil, err
	}

	// Determine the file type
	kind, err := filetype.Match(head)
	if err != nil {
		return nil, err
	}

	// If the file type is unknown, handle accordingly
	if kind == filetype.Unknown {
		return nil, errors.New("unknown file type")
	}

	return &kind, nil
}

func (h *StorageHandler) FilesList(c *fiber.Ctx) error {

	directory_id := uuid.FromStringOrNil(c.Params("directory_id"))
	if directory_id == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "directory ID is invalid"), nil)
	}

	files, err := h.IStorageInternal.FilesList(&directory_id)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if len(*files) == 0 {
		return SuccessResponse(c, NewStatus(200, "no files in this directory"), files)
	}

	return SuccessResponse(c, NewStatus(200), files)
}

func (h *StorageHandler) GetFileData(c *fiber.Ctx) error {

	file_id := uuid.FromStringOrNil(c.Params("file_id"))
	if file_id == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "file ID is invalid"), nil)
	}

	file, err := h.IStorageInternal.GetFileData(&file_id)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if file.ID == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "no such file"), nil)
	}

	return SuccessResponse(c, NewStatus(200), file)
}

func (h *StorageHandler) GetFile(c *fiber.Ctx) error {

	id := uuid.FromStringOrNil(c.Params("file_id"))

	fileModel, err := h.IStorageInternal.GetFileData(&id)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if fileModel == nil &&
		fileModel.ID == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	return c.SendFile(fileModel.FilePath(true))

}

func (h *StorageHandler) CreateFolder(c *fiber.Ctx) error {

	parent_id := uuid.FromStringOrNil(c.Params("parent_id"))
	folderModel := models.Directory{
		ParentID: parent_id,
	}.NewDirectory()

	err := c.BodyParser(&folderModel)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if folderModel == nil &&
		folderModel.ID == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	folderModel, err = h.IStorageInternal.CreateFolder(folderModel)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	return SuccessResponse(c, NewStatus(200, "Created"), folderModel)
}

func (h *StorageHandler) FoldersList(c *fiber.Ctx) error {

	parentID := uuid.FromStringOrNil(c.Params("parentID"))
	if parentID == uuid.Nil {
		return ErrorResponse(c, fiber.NewError(400, "directory ID is invalid"), nil)
	}

	folders, err := h.IStorageInternal.FoldersList(&parentID)
	if err != nil {
		return ErrorResponse(c, fiber.NewError(400, err.Error()), nil)
	}

	if len(*folders) == 0 {
		return SuccessResponse(c, NewStatus(200, "no folders in this directory"), folders)
	}

	return SuccessResponse(c, NewStatus(200), folders)
}
