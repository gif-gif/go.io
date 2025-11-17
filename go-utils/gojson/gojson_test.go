package gojson

import (
	"fmt"
	"testing"
)

type TestStruct struct {
	Field1 string  `json:"field1"`
	Field2 int64   `json:"field2"`
	Field3 int     `json:"field3"`
	Field4 bool    `json:"field4"`
	Field5 float32 `json:"field5"`
}

func TestGoJson(t *testing.T) {
	ts := TestStruct{
		Field1: "field1",
		Field2: 1234567890,
		Field3: 1234567890,
		Field4: true,
		Field5: 123.456,
	}
	data, err := Marshal(ts)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))
	prettyData, err := Marshal(ts, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(prettyData))

	// 反序列化
	// 反序列化到结构体
	// 反序列化到 map
	var m map[string]any
	err = Unmarshal(data, &m)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(m)
	ts2 := TestStruct{}
	err = Unmarshal(data, &ts2)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ts2)

	jsonStr, err := MarshalToString(ts2, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(jsonStr)

	ts3 := TestStruct{}
	err = UnmarshalFromString(jsonStr, &ts3)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ts3)
}
