package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_app/backend/errorz"
	"go_app/backend/helperz"
	"go_app/backend/models"
	"os"
	"reflect"
	"strings"

	"github.com/gofrs/uuid"
)

type IStorageInternal interface {
	ReadFilesList() ([]models.File, error)
	ReadFileData(id *uuid.UUID) (*models.File, error)
	UpdateFileData(id *uuid.UUID, fileModel *models.File) (*[]byte, error)
	DeleteFile(id *uuid.UUID) error

	CreateFolder(folder *models.Directory) (*models.Directory, error)
	ReadFoldersList() ([]models.Directory, error)
	GetFolders() (*[]models.Directory, error)
	ReadFolderData(id *uuid.UUID, name string) (*models.Directory, error)
	UpdateFolderData(id *uuid.UUID, fieldsToUpdate *models.Directory) error
	DeleteFolder(id *uuid.UUID) error
}

type StorageInternal struct{}

func NewStorage() IStorageInternal {
	return &StorageInternal{}
}

func (s *StorageInternal) ReadFilesList() ([]models.File, error) {

	jsonDB, err := os.ReadFile("C:/Users/ojpkm/Documents/go_app/Database/files.json")
	if err != nil {
		return nil, err
	}

	if jsonDB == nil {
		return nil, errors.New("file empty")
	}

	var file []models.File

	err = json.Unmarshal(jsonDB, &file)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *StorageInternal) ReadFoldersList() ([]models.Directory, error) {

	jsonDB, err := os.ReadFile("C:/Users/ojpkm/Documents/go_app/Database/directory.json")
	if err != nil {
		return nil, err
	}

	if jsonDB == nil {
		return nil, errors.New("file empty")
	}

	var folder []models.Directory

	err = json.Unmarshal(jsonDB, &folder)
	if err != nil {
		return nil, err
	}

	return folder, nil
}

func (s *StorageInternal) ReadFileData(id *uuid.UUID) (*models.File, error) {

	files, err := s.ReadFilesList()
	if err != nil {
		return nil, err
	}

	for _, v := range files {
		if v.ID == *id {
			return &v, nil
		}
	}

	return nil, errors.New("cannot find file")
}

func (s *StorageInternal) ReadFolderData(id *uuid.UUID, name string) (*models.Directory, error) {

	folders, err := s.ReadFoldersList()
	if err != nil {
		return nil, err
	}

	if name != "" {
		for _, v := range folders {
			if strings.ToLower(v.Name) == strings.ToLower(name) {
				return &v, nil
			}
		}
	}

	for _, v := range folders {
		if v.ID == *id {
			return &v, nil
		}
	}

	return nil, errors.New("cannot find folder")
}

func (s *StorageInternal) UpdateFolderData(id *uuid.UUID, fieldsToUpdate *models.Directory) error {

	folders, err := s.ReadFoldersList()
	if err != nil {
		return errorz.SendError(err)
	}

	var newStruct []byte

	for _, folder := range folders {
		if folder.ID == *id {

			go func(folder *models.Directory) error {

				newStruct, err = helperz.UpdateStruct(folder, fieldsToUpdate)
				if err != nil {
					return errorz.SendError(err)
				}

				fmt.Printf("reflect.TypeOf(&newStruct): %v\n", reflect.TypeOf(&newStruct))
				fmt.Printf("reflect.TypeOf(id): %v\n", reflect.TypeOf(id))

				err = helperz.DataBaseUpdate(&newStruct, "directory", id)
				if err != nil {
					return errorz.SendError(err)
				}

				return nil
			}(&folder)

		}

	}

	return nil
}

func (s *StorageInternal) UpdateFileData(id *uuid.UUID, fieldsToUpdate *models.File) (*[]byte, error) {

	files, err := s.ReadFilesList()
	if err != nil {
		return nil, errorz.SendError(err)
	}

	for i, v := range files {

		if v.ID == *id {

			bs, err := helperz.UpdateStruct(v, *fieldsToUpdate)
			if err != nil {
				return nil, errorz.SendError(err)
			}

			var updatedField models.File
			err = json.Unmarshal(bs, &updatedField)
			if err != nil {
				return nil, errorz.SendError(err)
			}

			files = append(files[:i], files[i+1:]...)
			files = append(files, updatedField)

			f, err := json.Marshal(files)
			if err != nil {
				return nil, errorz.SendError(err)
			}

			err = os.WriteFile("C:/Users/ojpkm/Documents/go_app/Database/files.json", f, os.ModeAppend)
			if err != nil {
				return nil, errorz.SendError(err)
			}

			if fieldsToUpdate.Name != "" {
				err := os.Rename(
					(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Storage/%s", v.ID)),
					(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Storage/%s", fieldsToUpdate.ID)))
				if err != nil {
					return nil, errorz.SendError(err)
				}

			}

			return &bs, nil
		}
	}

	return nil, errorz.SendError(errors.New("cannot update a file"))
}

func (s *StorageInternal) DeleteFile(id *uuid.UUID) error {

	files, err := s.ReadFilesList()
	if err != nil {
		return errorz.SendError(err)
	}

	for i, v := range files {

		if v.ID == *id {

			files = append(files[:i], files[i+1:]...)

			editedJSON, err := json.Marshal(files)
			if err != nil {
				return errorz.SendError(err)
			}

			err = os.WriteFile("C:/Users/ojpkm/Documents/go_app/Database/files.json", editedJSON, os.ModeAppend)
			if err != nil {
				return errorz.SendError(err)
			}

			err = os.Remove(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Storage/%s", v.ID))
			if err != nil {
				fmt.Printf("err: %v\n", err)
				return errorz.SendError(err)
			}

			return nil
		}
	}

	return errors.New("cannot find a file")
}

func (s *StorageInternal) DeleteFolder(id *uuid.UUID) error {

	folders, err := s.GetFolders()
	if err != nil {
		return errorz.SendError(err)
	}

	folds := *folders

	files, err := s.ReadFilesList()
	if err != nil {
		errorz.SendError(err)
	}

	for i, v := range folds {

		if v.ID == *id {

			folds = append(folds[:i], folds[i+1:]...)

			editedJSON, err := json.Marshal(folds)
			if err != nil {
				return errorz.SendError(err)
			}

			err = os.WriteFile("C:/Users/ojpkm/Documents/go_app/Database/directory.json", editedJSON, os.ModeAppend)
			if err != nil {
				return errorz.SendError(err)
			}

		}

	}

	for _, v := range files {

		if fileDirID := func() string {
			return id.String()
		}(); fileDirID == v.Directory {
			fmt.Printf("dirid: %v\n", fileDirID)
			err = s.DeleteFile(id)
			if err != nil {
				fmt.Printf("id: %v\n", id)
				return errorz.SendError(err)
			}
		}

	}

	return errors.New("cannot find a folder")
}

func (s *StorageInternal) CreateFolder(folder *models.Directory) (*models.Directory, error) {

	id, err := uuid.NewV7()
	if err != nil {
		return nil, errorz.SendError(err)
	}

	folder.ID = id

	fmt.Print("saved")
	helperz.PrettyPrint(folder)

	err = helperz.DataBaseInsert(folder, "directory")
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return folder, nil
}

func (s *StorageInternal) GetFolders() (*[]models.Directory, error) {

	folders, err := helperz.GetRecordsList("directory")
	if err != nil {
		return nil, errorz.SendError(err)
	}

	foldersModel := []models.Directory{}
	mainFound := false

	for _, v := range *folders {

		foldersModel = append(foldersModel, func() models.Directory {

			folder := models.Directory{}

			id, ok := v["id"].(string)
			if !ok {
				fmt.Errorf("cannot parse ID")
			}

			name, ok := v["name"].(string)
			if !ok {
				fmt.Errorf("cannot parse ID")
			}

			folder.ID = uuid.FromStringOrNil(id)
			folder.Name = name

			return folder
		}())

		if strings.ToLower(v["name"].(string)) == "main" {
			mainFound = true
		}

	}

	if !mainFound {

		main, err := s.CreateFolder(func() *models.Directory {
			return &models.Directory{
				Name: "Main",
			}
		}())
		if err != nil {
			fmt.Errorf("cannot create Main")
		}

		foldersModel = append(foldersModel, *main)

	}

	return &foldersModel, nil
}
