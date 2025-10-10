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

// 在项目入口统一配置
var JSON = jsoniter.ConfigCompatibleWithStandardLibrary

func Pretty(strJson string) string {
	var out bytes.Buffer
	if err := json.Indent(&out, ([]byte)(strJson), "", "    "); err != nil {
		return ""
	}
	return out.String()
}

func Marshal(obj any, pretty bool) string {
	var data []byte
	if pretty {
		data, _ = JSON.MarshalIndent(obj, "", "    ")
	} else {
		data, _ = JSON.Marshal(obj)
	}
	return string(data)
}

func UnmarshalFromFile(filePath string, val any) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	return JSON.Unmarshal(data, val)
}

func Unmarshal(data string, val any) error {
	return JSON.UnmarshalFromString(data, val)
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
