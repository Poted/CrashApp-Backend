package internal

import (
	"go_app/backend/db"
	"go_app/backend/errorz"
	"go_app/backend/models"
	"os"

	"github.com/gofrs/uuid"
)

type IStorageInternal interface {

	// Files

	SaveFile(fileModel *models.File, saveFileFunc func(path string) error) (*models.File, error)
	FilesList(directoryID *uuid.UUID) (*[]models.File, error)
	GetFileData(fileID *uuid.UUID) (*models.File, error)

	// Folders

	CreateFolder(folderModel *models.Directory) (*models.Directory, error)
	FoldersList(parentID *uuid.UUID) (*[]models.Directory, error)

	// ReadFilesList() (*[]models.File, error)
	// ReadFileData(id *uuid.UUID) (*models.File, error)
	// UpdateFileData(id *uuid.UUID, fileModel *models.File) (*[]byte, error)
	// DeleteFile(id *uuid.UUID) error

	// CreateFolder(folder *models.Directory) (*models.Directory, error)
	// ReadFoldersList() (*[]models.Directory, error)
	// ReadFolderData(id *uuid.UUID, name string) (*models.Directory, error)
	// UpdateFolderData(id *uuid.UUID, fieldsToUpdate *models.Directory) error
	// DeleteFolder(id *uuid.UUID) error
}

type StorageInternal struct{}

func NewStorage() IStorageInternal {
	return &StorageInternal{}
}

func (s *StorageInternal) SaveFile(fileModel *models.File, saveFileFunc func(path string) error) (*models.File, error) {

	tx := db.Database.Begin()

	// Save file data to database
	err := tx.
		Save(&fileModel).
		Error
	if err != nil {
		tx.Rollback()
		return nil, errorz.SendError(err)
	}

	// Create directories to store file
	err = os.MkdirAll(fileModel.FilePath(false), 0755)
	if err != nil {
		tx.Rollback()
		return nil, errorz.SendError(err)
	}

	// Save file to storage folder
	err = saveFileFunc(fileModel.FilePath(true))
	if err != nil {
		tx.Rollback()
		return nil, errorz.SendError(err)
	}

	err = tx.Commit().Error

	return fileModel, err
}

func (s *StorageInternal) FilesList(directoryID *uuid.UUID) (*[]models.File, error) {

	files := []models.File{}

	err := db.Database.
		Model(&models.File{}).
		Where("directory_id = ?", directoryID).
		Scan(&files).
		Error

	return &files, err
}

func (s *StorageInternal) GetFileData(fileID *uuid.UUID) (*models.File, error) {

	file := models.File{}

	err := db.Database.Debug().
		Model(&models.File{}).
		Where("id = ?", fileID).
		Scan(&file).
		Error

	return &file, err
}

func (s *StorageInternal) CreateFolder(folderModel *models.Directory) (*models.Directory, error) {

	err := db.Database.
		Save(&folderModel).
		Error

	return folderModel, err
}

func (s *StorageInternal) FoldersList(parentID *uuid.UUID) (*[]models.Directory, error) {

	folders := []models.Directory{}

	err := db.Database.
		Model(&models.File{}).
		Where("parent_id = ?", parentID).
		Scan(&folders).
		Error

	return &folders, err
}
