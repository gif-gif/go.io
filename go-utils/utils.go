package goutils

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"math"
	"math/rand"
	"strings"
	"time"
)

// 通用三目运算
func IfNot[T any](isTrue bool, a, b T) T {
	if isTrue {
		return a
	}
	return b
}

func IfString(isTrue bool, a, b string) string {
	if isTrue {
		return a
	}
	return b
}

func IfInt(isTrue bool, a, b int) int {
	if isTrue {
		return a
	}
	return b
}

func IfFloat32(isTrue bool, a, b float32) float32 {
	if isTrue {
		return a
	}
	return b
}

func IfFloat64(isTrue bool, a, b float64) float64 {
	if isTrue {
		return a
	}
	return b
}

func ReverseArray(arr []*interface{}) {
	for i, j := 0, len(arr)-1; i <= j; i, j = i+1, j-1 {
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func PadStart(str, pad string, length int) string {
	if len(str) >= length {
		return str
	}
	return strings.Repeat(pad, length-len(str)) + str
}

func MinInt64(a, b int64) int64 {
	return gconv.Int64(math.Min(gconv.Float64(a), gconv.Float64(b)))
}

func MinInt(a, b int) int {
	return gconv.Int(math.Min(gconv.Float64(a), gconv.Float64(b)))
}

func MaxInt64(a, b int64) int64 {
	return gconv.Int64(math.Max(gconv.Float64(a), gconv.Float64(b)))
}

func MaxInt(a, b int) int {
	return gconv.Int(math.Max(gconv.Float64(a), gconv.Float64(b)))
}

func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}
