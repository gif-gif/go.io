package goutils

import (
	"bytes"
	"fmt"
	"slices"
	"sort"

	"github.com/samber/lo"
)

func CloneMap[TKey comparable, TValue any](m1 map[TKey]TValue) map[TKey]TValue {
	m2 := make(map[TKey]TValue, len(m1))
	for k, v := range m1 {
		m2[k] = v
	}
	return m2
}

func ArrayJoin[T any](arr []T, sep string) string {
	buf := bytes.NewBuffer(nil)

	for idx, val := range arr {
		if idx == 0 {
			buf.WriteString(fmt.Sprintf("%v", val))
		} else {
			buf.WriteString(fmt.Sprintf("%v%v", sep, val))
		}
	}

	return buf.String()
}

// 元素都转换成字符串比较
func IsInArray[T comparable](arr []T, target T) bool {
	return slices.Contains(arr, target)
}

// 条件满足任意元素 exists func(target T) bool 返回true时返回true
//
// 适合判断数组中存储复杂对象，判断条件定义情况
//
// 用以下代替
//
//	slices.ContainsFunc(arr, func(t T) bool {
//
//	})
func IsInArrayX[T any](arr []T, exists func(target T) bool) bool {
	for _, t := range arr {
		if exists(t) {
			return true
		}
	}
	return false
}

// 条件满足任意元素 exists func(target *T) bool 返回true时返回true
//
// 适合判断数组中存储复杂对象，判断条件定义情况,数组元素是指针类型时用
//
// 用以下代替
//
//	slices.ContainsFunc(arr, func(t T) bool {
//
//	})
func IsInArrayXX[T any](arr []*T, exists func(target *T) bool) bool {
	for _, t := range arr {
		if exists(t) {
			return true
		}
	}
	return false
}

func ReverseArray(arr []*interface{}) {
	for i, j := 0, len(arr)-1; i <= j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

// 插入排序函数，使用泛型指定元素类型
func InsertionSort[T any](arr []T, less func(T, T) bool) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1

		// 将比key大的元素向后移动一位
		for j >= 0 && less(arr[j], key) == false {
			arr[j+1] = arr[j]
			j--
		}

		// 插入关键元素到正确的位置
		arr[j+1] = key
	}
}

// 定义一个泛型排序函数
func GenericSort[T any](arr []T, less func(T, T) bool) {
	sort.Slice(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
}

// 两个数组是否相等，判断长度一样的两个数组 元素是否完全相同，顺序可以不同
func IsEqualArray[T comparable](arr1 []T, arr2 []T) bool {
	if len(arr1) != len(arr2) {
		return false
	}

	for _, v := range arr1 {
		if !lo.Contains(arr2, v) {
			return false
		}
	}
	return true
}

func RemoveArrayDuplicateValue[T any](arr []T) []T {
	if len(arr) <= 1 {
		return arr
	}
	seen := make(map[string]struct{}, len(arr))
	result := make([]T, 0, len(arr))

	for _, val := range arr {
		key := fmt.Sprintf("%v", val)
		if _, ok := seen[key]; !ok {
			result = append(result, val)
			seen[key] = struct{}{}
		}
	}
	return result
}

func SplitStringArray(arr []string, size int) (list [][]string) {
	l := len(arr)

	if l == 0 {
		list = make([][]string, 0)
		return
	}

	if l < size {
		list = [][]string{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitIntArray(arr []int, size int) (list [][]int) {
	l := len(arr)

	if l == 0 {
		list = make([][]int, 0)
		return
	}

	if l < size {
		list = [][]int{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitInt64Array(arr []int64, size int) (list [][]int64) {
	l := len(arr)

	if l == 0 {
		list = make([][]int64, 0)
		return
	}

	if l < size {
		list = [][]int64{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

func SplitArray(arr []interface{}, size int) (list [][]interface{}) {
	l := len(arr)

	if l == 0 {
		list = make([][]interface{}, 0)
		return
	}

	if l < size {
		list = [][]interface{}{arr}
		return
	}

	var (
		offset int
	)

	for {
		if offset+size >= l {
			list = append(list, arr[offset:])
			break
		}

		list = append(list, arr[offset:offset+size])

		offset += size
	}

	return
}

//
//// Insert 在指定位置插入元素
//// 支持负数索引，-1 表示最后一个位置，-2 表示倒数第二个位置
//func InsertArray[T any](slice []T, index int, element T) []T {
//	// 处理负数索引
//	if index < 0 {
//		index = len(slice) + index + 1
//	}
//
//	// 边界检查
//	if index < 0 {
//		index = 0
//	} else if index > len(slice) {
//		index = len(slice)
//	}
//
//	// 插入元素
//	return append(slice[:index], append([]T{element}, slice[index:]...)...)
//}

// Comparator 函数类型，接受两个参数并返回整数：
// - 负数表示 a < b
// - 零表示 a == b
// - 正数表示 a > b
type Comparator[T any] func(a, b T) int

// SortWith 使用多个比较函数对切片进行排序
func SortWith[T any](fns []Comparator[T], list []T) []T {
	result := make([]T, len(list))
	copy(result, list)

	// 使用自定义排序逻辑
	slices.SortStableFunc(result, func(a, b T) int {
		for _, fn := range fns {
			if cmp := fn(a, b); cmp != 0 {
				return cmp
			}
		}
		return 0 // 所有比较函数都返回 0，保持原顺序
	})

	return result
}
