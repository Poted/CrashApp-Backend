package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"go_app/backend/errorz"
	"go_app/backend/helperz"
	"go_app/backend/models"
	"os"

	"github.com/google/uuid"
)

type IStorageInternal interface {
	ReadFilesList() ([]models.File, error)
	ReadFileData(id *uuid.UUID) (*models.File, error)
	UpdateFileData(id *uuid.UUID, fileModel *models.File) (*[]byte, error)
	DeleteFile(id *uuid.UUID) error

	ReadFoldersList() ([]models.Directory, error)
	// DataBaseInsert(structModel interface{}, tableName string) error
}

type StorageInternal struct{}

func New() IStorageInternal {
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

	jsonDB, err := os.ReadFile("C:/Users/ojpkm/Documents/go_app/Database/directories.json")
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

// func (s *StorageInternal) DataBaseInsert(structModel interface{}, tableName string) error {

// 	jsonDB, err := os.OpenFile(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Database/%s.json", tableName), os.O_RDWR|os.O_CREATE, 0644)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer jsonDB.Close()

// 	// Decode the existing JSON data from the file
// 	var files []interface{}
// 	decoder := json.NewDecoder(jsonDB)
// 	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
// 		log.Fatal(err)
// 	}

// 	// Append the new record to the existing data
// 	files = append(files, structModel)

// 	// Seek to the beginning of the file to overwrite the existing data
// 	if _, err := jsonDB.Seek(0, 0); err != nil {
// 		log.Fatal(err)
// 	}

// 	// Create a JSON encoder and encode the updated data to the file
// 	encoder := json.NewEncoder(jsonDB)
// 	if err := encoder.Encode(files); err != nil {
// 		log.Fatal(err)
// 	}

// 	return nil
// }