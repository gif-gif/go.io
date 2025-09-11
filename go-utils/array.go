package goutils

import "slices"

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
