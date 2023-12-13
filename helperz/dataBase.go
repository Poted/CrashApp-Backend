package helperz

import (
	"encoding/json"
	"fmt"
	"go_app/backend/errorz"
	"log"
	"os"
	"path/filepath"

	"github.com/gofrs/uuid"
)

var dbPath = "C:/Users/ojpkm/Documents/go_app/Database"

// Remember to 'defer table.Close()';
func getTable(tableName string) (*os.File, error) {

	path := dbPath + "/%s.json"

	jsonDB, err := os.OpenFile(fmt.Sprintf(path, tableName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return jsonDB, nil
}

func GetRecordsList(tableName string) (*[]map[string]interface{}, error) {

	table, err := getTable(tableName)
	if err != nil {
		return nil, errorz.SendError(err)
	}
	defer table.Close()

	decoder := json.NewDecoder(table)

	var Ifiles []interface{}
	var files []map[string]interface{}

	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		return nil, errorz.SendError(err)
	}

	for _, file := range Ifiles {

		f, ok := file.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("cannot read data")
		}
		files = append(files, f)
	}

	return &files, nil
}

func GetSingleRecord(ID *uuid.UUID, tableName string) (*map[string]interface{}, error) {

	records, err := GetRecordsList(tableName)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	for _, r := range *records {

		for key, value := range r {
			if key == "id" {
				if value == ID {
					return &r, nil
				}
			}
		}

	}

	return nil, fmt.Errorf("not found")
}

func DataBaseInsert(structModel interface{}, tableName string) error {

	table, err := getTable(tableName)
	if err != nil {
		return errorz.SendError(err)
	}
	defer table.Close()

	// Decode the existing JSON data from the file
	var files []interface{}
	decoder := json.NewDecoder(table)

	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		return errorz.SendError(err)
	} // stond bugi wyszli

	// Append the new record to the existing data
	files = append(files, structModel)

	// Seek to the beginning of the file to overwrite the existing data
	if _, err := table.Seek(0, 0); err != nil {
		return errorz.SendError(err)
	}

	// Create a JSON encoder and encode the updated data to the file
	encoder := json.NewEncoder(table)
	if err := encoder.Encode(files); err != nil {
		return errorz.SendError(err)
	}

	return nil
}

func DataBaseUpdate(structToSave any, tableName string, id *uuid.UUID) error {

	table, err := getTable(tableName)
	if err != nil {
		return errorz.SendError(err)
	}
	defer table.Close()

	// Decode the existing JSON data from the file
	var files []map[string]interface{}
	decoder := json.NewDecoder(table)
	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	for _, v := range files {
		if v["id"] == id.String() {

			fmt.Printf("v: %v\n", v)

			newStruct, err := UpdateStruct(v, structToSave)
			if err != nil {
				return errorz.SendError(err)
			}

			err = json.Unmarshal(newStruct, &v)
			if err != nil {
				return errorz.SendError(err)
			}

		}
	}

	// Seek to the beginning of the file to overwrite the existing data
	if _, err := table.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Create a JSON encoder and encode the updated data to the file
	encoder := json.NewEncoder(table)
	if err := encoder.Encode(files); err != nil {
		log.Fatal(err)
	}

	return nil
}

func GetStorage(tableName string) error {

	storagePath := "C:/Users/ojpkm/Documents/go_app/Storage"

	err := filepath.WalkDir(storagePath, visitFile)
	if err != nil {
		return errorz.SendError(err)
	}

	return nil
}

func visitFile(fp string, fi os.DirEntry, err error) error {

	if err != nil {
		fmt.Println(err) // can't walk here,
		return nil       // but continue walking elsewhere
	}

	if fi.IsDir() {
		return nil // not a file. ignore.
	}

	return nil
}
