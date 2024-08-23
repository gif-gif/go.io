package main

import "fmt"

func main() {
	testT()
}

// 泛型类型别名预览示例
// 使用泛型类型别名
func testT() {
	// 定义一个泛型类型别名，表示任意类型T的指针
	type Ptr[T any] *T

	// 定义一个泛型类型别名，表示任意类型T的切片
	type Slice[T any] []T
	// 创建一个int类型的指针
	var p Ptr[int] = new(int)
	*p = 42

	// 创建一个string类型的切片
	var s Slice[string] = []string{"Hello", "World"}

	fmt.Println(*p) // 输出: 42
	fmt.Println(s)  // 输出: [Hello World]
}
