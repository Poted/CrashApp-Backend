package helperz

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
)

type IDataBase interface {
	StructIterate(stc any) (map[string]interface{}, error)
	DataBaseInsert(structModel interface{}, tableName string) error
}

func DataBaseInsert(structModel interface{}, tableName string) error {

	jsonDB, err := os.OpenFile(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Database/%s.json", tableName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonDB.Close()

	// Decode the existing JSON data from the file
	var files []interface{}
	decoder := json.NewDecoder(jsonDB)
	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	// Append the new record to the existing data
	files = append(files, structModel)

	// Seek to the beginning of the file to overwrite the existing data
	if _, err := jsonDB.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Create a JSON encoder and encode the updated data to the file
	encoder := json.NewEncoder(jsonDB)
	if err := encoder.Encode(files); err != nil {
		log.Fatal(err)
	}

	return nil
}

func DataBaseUpdate(structToSave interface{}, tableName string, id *uuid.UUID) error {

	jsonDB, err := os.OpenFile(fmt.Sprintf("C:/Users/ojpkm/Documents/go_app/Database/%s.json", tableName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer jsonDB.Close()

	// Decode the existing JSON data from the file
	// var files []interface{}
	var files []map[string]interface{}
	decoder := json.NewDecoder(jsonDB)
	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		log.Fatal(err)
	}

	// Append the new record to the existing data
	// files = append(files, structToSave)

	// recordsMap, err := StructIterate(files)

	// recordsMap["id"] = 2

	for _, v := range files {
		fmt.Printf("v: %v\n", v)
		if v["id"] == id {
			fmt.Printf("v: %v\n", v)
		}
	}

	// Seek to the beginning of the file to overwrite the existing data
	if _, err := jsonDB.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	// Create a JSON encoder and encode the updated data to the file
	encoder := json.NewEncoder(jsonDB)
	if err := encoder.Encode(files); err != nil {
		log.Fatal(err)
	}

	return nil
}
