package goio

import (
	"encoding/json"
	"fmt"
	goutils "github.com/gif-gif/go.io/go-utils"
	"github.com/go-playground/validator/v10"

	"strings"
)

func ValidationMessage(err error, messages map[string]string) string {
	if v, ok := err.(*json.UnmarshalTypeError); ok {
		return fmt.Sprintf("请求参数 %s 的类型是 %s, 不是 %s", v.Field, v.Type, v.Value)
	}

	if v, ok := err.(validator.ValidationErrors); ok {
		for _, i := range v {
			field := goutils.Camel2Case(i.Field())
			key := fmt.Sprintf("%s_%s", field, strings.ToLower(i.Tag()))
			if msg, ok := messages[key]; ok {
				return msg
			}
			return fmt.Sprintf("%s %s", field, i.Tag())
		}
	}

	return err.Error()
}
