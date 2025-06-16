package goconvert

import (
	"encoding/json"
)

func ConvertToMap(data interface{}) map[string]interface{} {
	var mapResult map[string]interface{}
	dataBytes, _ := json.Marshal(data)
	_ = json.Unmarshal(dataBytes, &mapResult)
	return mapResult
}

func ConvertMapToStruct[T any](data map[string]interface{}) (*T, error) {
	var t T
	fieldsBytes, err := json.Marshal(data)
	if err != nil {
		return &t, err
	}

	err = json.Unmarshal(fieldsBytes, &t)
	if err != nil {
		return &t, err
	}
	return &t, nil
}
