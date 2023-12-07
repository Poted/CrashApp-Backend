package helperz

import (
	"encoding/json"
	"go_app/backend/errorz"
	"reflect"
)

type IHelperz interface {
	StructIterate(stc any) (map[string]interface{}, error)
}

// Converting a struct into a map[string]interface representation of a json.
func StructIterate(stc any) (map[string]interface{}, error) {

	var stcMap map[string]interface{}
	stcJSON, err := json.Marshal(&stc)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	err = json.Unmarshal(stcJSON, &stcMap)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return stcMap, nil
}

// Updating a struct with new values of whole struct or its separate values.
func UpdateStruct(legacyStruct any, updateFields any) ([]byte, error) {

	// Converting struct into map[string]interface{}
	mapLS, err := StructIterate(legacyStruct)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	// Converting new fields into map[string]interface{}
	mapUF, err := StructIterate(updateFields)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	// Updating struct; Switching old values with their new value.
	for key, value := range mapUF {

		if !reflect.ValueOf(value).IsZero() {
			mapLS[key] = value
		}

	}

	// Returning as a []byte for later converting to legacy struct.
	jsn, err := json.Marshal(mapLS)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return jsn, nil
}
