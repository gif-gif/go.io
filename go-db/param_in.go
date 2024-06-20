package godb

import (
	"strings"
)

func GenerateSliceIn[T any](srcItems []T) (string, []any) {
	if len(srcItems) == 0 {
		return "", nil
	}

	targetItems := make([]any, 0, len(srcItems))
	builder := strings.Builder{}
	builder.WriteString(" ( ")
	for _, item := range srcItems {
		builder.WriteString("?,")
		targetItems = append(targetItems, item)
	}

	targetString := builder.String()
	targetString = targetString[:len(targetString)-1]

	return targetString + " ) ", targetItems
}

func GenerateSliceInEx[T any](fieldName string, srcItems []T) (string, []any) {
	if len(srcItems) == 0 {
		return "", nil
	}

	targetItems := make([]any, 0, len(srcItems))
	builder := strings.Builder{}
	builder.WriteString(" ( ")
	for _, item := range srcItems {
		builder.WriteString("?,")
		targetItems = append(targetItems, item)
	}

	targetString := builder.String()
	targetString = targetString[:len(targetString)-1]

	return targetString + " ) ", targetItems
}
