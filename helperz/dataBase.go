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
func getTable(tableName string) *os.File {

	dbPath = dbPath + "/%s.json"

	jsonDB, err := os.OpenFile(fmt.Sprintf(dbPath, tableName), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}

	return jsonDB
}

func DataBaseInsert(structModel interface{}, tableName string) error {

	table := getTable(tableName)
	defer table.Close()

	// Decode the existing JSON data from the file
	var files []interface{}
	decoder := json.NewDecoder(table)
	if err := decoder.Decode(&files); err != nil && err.Error() != "EOF" {
		return errorz.SendError(err)
	}

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

	table := getTable(tableName)
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

func RestoreStorage(tableName string) error {

	storagePath := fmt.Sprint("C:/Users/ojpkm/Documents/go_app/Storage")

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

	// sp, err := fi.Info()
	// if err != nil {
	// 	return errorz.SendError(err)
	// }

	// name := sp.Name()
	// size := sp.Size()

	// fmt.Printf("fp: %v\n", fp)

	// fmt.Printf("name: %v\n", name)
	// fmt.Printf("size: %v\n", size)

	// file := models.File{
	// 	ID:        uuid.FromStringOrNil(name),
	// 	Name:      fmt.Sprint(name + "restored"),
	// 	Size:      size,
	// 	Directory: "",
	// }

	// fmt.Printf("file: %v\n", file)

	// DataBaseInsert(file, "files-restored")

	return nil
}
