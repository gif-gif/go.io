package gobean

import (
	"encoding/json"
	"fmt"
	"golang.org/x/exp/constraints"
)

const (
	OpEquals         = "equals"
	OpLessOrEquals   = "lessEqual"
	OpGreaterOrEqual = "greaterEqual"
	OpGreater        = "greater"
	OpLess           = "less"
)

type NumberRange[T any] struct {
	Min *T `json:"min,optional"`
	Max *T `json:"max,optional"`
}

type Int64Range struct {
	Min *int64 `json:"min,optional"`
	Max *int64 `json:"max,optional"`
}

type Float64Range struct {
	Min *float64 `json:"min,optional"`
	Max *float64 `json:"max,optional"`
}

type Int64Method struct {
	Method string `json:"method,optional"`
	Value  int64  `json:"value,optional"`
}

func (r *Int64Method) Check(val int64) bool {
	return CheckValue(r, val)
}
func (r *Int64Method) GetValue() int64 {
	return r.Value
}

func (r *Int64Method) GetMethod() string {
	return r.Method
}

type IOperationMethod[T any] interface {
	GetMethod() string
	GetValue() T
}

// 大于，大于等于，小于，小于等于，等于
func CheckValue(r IOperationMethod[int64], val int64) bool {
	method := r.GetMethod()
	value := r.GetValue()
	switch method {
	case OpEquals:
		return val == value
	case OpLessOrEquals:
		return val <= value
	case OpGreaterOrEqual:
		return val >= value
	case OpGreater:
		return val > value
	case OpLess:
		return val < value
	}
	return false
}

// IsInRange 检查给定值是否在指定的范围内
//
// 如果 Min 为 nil，则没有下限限制
//
// 如果 Max 为 nil，则没有上限限制
//
// 如果 Min 和 Max 都为 nil，则认为没有限制，任何值都符合条件
func IsInRange[T constraints.Ordered](value T, rng NumberRange[T]) bool {
	// 如果 Min 和 Max 都为 nil，则没有限制
	if rng.Min == nil && rng.Max == nil {
		return true
	}

	// 检查下限
	if rng.Min != nil && value < *rng.Min {
		return false
	}

	// 检查上限
	if rng.Max != nil && value > *rng.Max {
		return false
	}

	return true
}

// CheckRanges 检查给定值是否满足所有范围条件
//
// 有符合条件的返回true，否则false
func CheckRanges[T constraints.Ordered](value T, ranges []NumberRange[T]) bool {
	// 如果没有条件，则视为没有限制
	if len(ranges) == 0 {
		return false
	}

	// 检查是否满足所有条件
	for _, rng := range ranges {
		if IsInRange(value, rng) {
			return true
		}
	}

	return false
}

// 创建一个帮助函数，用于更方便地创建 NumberRange
func NewNumberRange[T constraints.Ordered](min, max *T) NumberRange[T] {
	return NumberRange[T]{
		Min: min,
		Max: max,
	}
}

// 有符合条件的返回true，否则false
func CheckValueAgainstRanges[T constraints.Ordered](value T, jsonRanges string) (bool, error) {
	// 如果 JSON 字符串为空，视为没有限制
	if jsonRanges == "" {
		return false, nil
	}

	// 解析 JSON 数组
	var ranges []NumberRange[T]
	err := json.Unmarshal([]byte(jsonRanges), &ranges)
	if err != nil {
		return false, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	// 空数组视为没有限制
	if len(ranges) == 0 {
		return false, nil
	}

	// 检查是否满足所有条件
	for _, rng := range ranges {
		if IsInRange(value, rng) {
			return true, nil
		}
	}

	return false, nil
}
