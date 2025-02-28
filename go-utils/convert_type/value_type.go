package goconvert

import (
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"strings"
)

const (
	ValueTypeInt     = "int"
	ValueTypeString  = "string"
	ValueTypeText    = "text"
	ValueTypeStrings = "string[]"
	ValueTypeBool    = "bool"
	ValueTypeDecimal = "decimal"
	ValueTypeMap     = "map"
	ValueTypeJson    = "json"
	ValueTypeByte    = "byte"
	ValueTypeBytes   = "bytes"
	ValueTypeSelect  = "select"
	ValueTypeUnknown = "unknown"
)

// ConvertValue converts a string value to a specific type based on the provided valueType
func ConvertValue(value string, valueType string, result ...interface{}) (interface{}, error) {
	switch valueType {
	case ValueTypeInt:
		return gconv.Int64(value), nil
	case ValueTypeString:
	case ValueTypeSelect:
	case ValueTypeText:
		return value, nil
	case ValueTypeDecimal:
		return gconv.Float64(value), nil
	case ValueTypeBool:
		return gconv.Bool(value), nil
	case ValueTypeStrings:
		return ConvertToStrings(value)
	case ValueTypeByte:
		return gconv.Byte(value), nil
	case ValueTypeBytes:
		return gconv.Bytes(value), nil
	case ValueTypeMap:
		return convertToMap(value)
	case ValueTypeJson:
		if len(result) == 0 {
			return convertToJson(value)
		}
		return ConvertToJson(value, result[0])
	case ValueTypeUnknown:
		return value, nil
	default:
		if strings.HasPrefix(valueType, ValueTypeSelect) {
			return value, nil
		}
		return nil, fmt.Errorf("unsupported value type: %s", valueType)
	}

	return nil, nil
}

// Helper functions for converting string to specific types

func ConvertToStrings(value string) ([]string, error) {
	var result []string
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, fmt.Errorf("cannot convert string to []string: %w", err)
	}
	return result, nil
}

func convertToMap(value string) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, fmt.Errorf("cannot convert string to map: %w", err)
	}
	return result, nil
}

func convertToJson(value string) (interface{}, error) {
	var result interface{}
	return ConvertToJson(value, result)
}

func ConvertToJson(value string, result interface{}) (interface{}, error) {
	err := json.Unmarshal([]byte(value), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return result, nil
}

// Example usage
func main() {
	// 测试一些字符串转换的例子
	intStr := "42"
	floatStr := "3.14"
	stringVal := "hello"
	mapStr := `{"name":"John","age":30}`
	jsonStr := `{"name":"John","age":30,"address":{"city":"New York","zip":"10001"}}`

	intResult, err := ConvertValue(intStr, ValueTypeInt)
	fmt.Printf("Int64: %v, %T, Error: %v\n", intResult, intResult, err)

	floatResult, err := ConvertValue(floatStr, ValueTypeDecimal)
	fmt.Printf("Float32: %v, %T, Error: %v\n", floatResult, floatResult, err)

	stringResult, err := ConvertValue(stringVal, ValueTypeString)
	fmt.Printf("String: %v, %T, Error: %v\n", stringResult, stringResult, err)

	byteResult, err := ConvertValue(intStr, ValueTypeByte)
	fmt.Printf("Byte: %v, %T, Error: %v\n", byteResult, byteResult, err)

	bytesResult, err := ConvertValue(stringVal, ValueTypeBytes)
	fmt.Printf("Bytes: %v, %T, Error: %v\n", bytesResult, bytesResult, err)

	mapResult, err := ConvertValue(mapStr, ValueTypeMap)
	fmt.Printf("Map: %v, %T, Error: %v\n", mapResult, mapResult, err)

	jsonResult, err := ConvertValue(jsonStr, ValueTypeJson)
	fmt.Printf("JSON: %v, %T, Error: %v\n", jsonResult, jsonResult, err)
}
