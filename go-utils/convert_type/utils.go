package goconvert

import (
	"encoding/json"

	"github.com/gif-gif/go.io/go-utils/gojson"
	"github.com/mitchellh/mapstructure"
)

func ConvertToMap(data interface{}) map[string]interface{} {
	var mapResult map[string]interface{}
	dataBytes, _ := json.Marshal(data)
	_ = gojson.Unmarshal(dataBytes, &mapResult)
	return mapResult
}

// for json tag 慢
func ConvertMapToStruct[T any](data map[string]interface{}) (*T, error) {
	var t T
	fieldsBytes, err := gojson.Marshal(data)
	if err != nil {
		return &t, err
	}

	err = gojson.Unmarshal(fieldsBytes, &t)
	if err != nil {
		return &t, err
	}
	return &t, nil
}

// for mapstructure tag 快
func ConvertMapToStructEx[T any](data map[string]interface{}) (*T, error) {
	var t T
	// 转换为结构体
	err := mapstructure.Decode(data, &t)
	if err != nil {
		return &t, err
	}
	return &t, nil
}
