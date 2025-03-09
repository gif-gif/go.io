package goconvert

import "encoding/json"

func ConvertToMap(data interface{}) map[string]interface{} {
	var mapResult map[string]interface{}
	dataBytes, _ := json.Marshal(data)
	_ = json.Unmarshal(dataBytes, &mapResult)
	return mapResult
}
