package gojson

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/oliveagle/jsonpath"
)

var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func Pretty(strJson string) string {
	var out bytes.Buffer
	if err := json.Indent(&out, ([]byte)(strJson), "", "    "); err != nil {
		return ""
	}
	return out.String()
}

func MarshalToString(obj any, pretty ...bool) (string, error) {
	var data []byte
	isPretty := len(pretty) > 0 && pretty[0]
	var err error
	if isPretty {
		data, err = JSON.MarshalIndent(obj, "", "    ")
	} else {
		data, err = JSON.Marshal(obj)
	}
	return string(data), err
}

func Marshal(obj any, pretty ...bool) ([]byte, error) {
	var data []byte
	isPretty := len(pretty) > 0 && pretty[0]
	var err error
	if isPretty {
		data, err = JSON.MarshalIndent(obj, "", "    ")
	} else {
		data, err = JSON.Marshal(obj)
	}
	return data, err
}

func UnmarshalFromFile(filePath string, val any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return JSON.Unmarshal(data, val)
}

func UnmarshalFromString(data string, val any) error {
	return JSON.UnmarshalFromString(data, val)
}

func Unmarshal(data []byte, val any) error {
	return JSON.Unmarshal(data, val)
}

func UnmarshalFromGzip(data []byte, val any) error {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer reader.Close()

	orginData, err := io.ReadAll(reader)
	if err != nil {
		return err
	}

	return JSON.Unmarshal(orginData, val)
}

func PathLookup[T any](obj interface{}, jpath string) (T, error) {
	var tval T
	var err error
	var ok bool

	val, err := jsonpath.JsonPathLookup(obj, jpath)
	if err != nil {
		return tval, err
	}

	tval, ok = val.(T)
	if !ok {
		return tval, errors.New("type conversion failed")
	}

	return tval, nil
}
