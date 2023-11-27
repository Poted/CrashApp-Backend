package helperz

import (
	"encoding/json"
	"go_app/backend/errorz"
	"reflect"
)

type IHelperz interface {
	StructIterate(stc any) (map[string]interface{}, error)
}

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

func UpdateStruct(legacyStruct any, updateFields any) ([]byte, error) {

	mapLS, err := StructIterate(legacyStruct)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	mapUF, err := StructIterate(updateFields)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	for key, value := range mapUF {

		if !reflect.ValueOf(value).IsZero() {
			mapLS[key] = value
		}

	}

	jsn, err := json.Marshal(mapLS)
	if err != nil {
		return nil, errorz.SendError(err)
	}

	return jsn, nil
}
