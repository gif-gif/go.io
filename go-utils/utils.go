package goutils

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"reflect"
	"runtime"
	"sort"
	"strings"
	"time"
)

// BcryptHash 使用 bcrypt 对密码进行加密
func BcryptHash(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

// BcryptCheck 对比明文密码和数据库的哈希值
func BcryptCheck(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// 常用签名验证, sign md5 小写
func CheckSign(secret string, linkSignTimeout int64, ts int64, sign string) bool {
	if linkSignTimeout == 0 {
		linkSignTimeout = 20
	}
	tsStep := time.Now().Unix() - ts
	if math.Abs(gconv.Float64(tsStep)) > gconv.Float64(linkSignTimeout) { //连接失效
		return false
	}
	serverSign := Md5([]byte(gconv.String(ts) + secret))
	return serverSign == sign
}

// 元素都转换成字符串比较
func IsInArray[T any](arr []T, target T) bool {
	for _, t := range arr {
		tt := gconv.String(target)
		if gconv.String(t) == tt {
			return true
		}
	}
	return false
}

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

// 通过反射获取结构体字段的值
func GetFieldValue(config interface{}, fieldName string) (interface{}, error) {
	v := reflect.ValueOf(config)

	// 确保传入的是一个指针
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// 确保传入的是结构体
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("expected a struct, but got %s", v.Kind())
	}

	// 获取字段值
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		return nil, fmt.Errorf("no such field: %s in struct", fieldName)
	}

	return field.Interface(), nil
}

// 使用反射复制A结构到B结构，前提是两个结构体字段数量和类型完全相同
// 如：使用反射将 VO 转换为 DTO
// 反射(reflect)虽爽，但很贵,性能会有损失
func CopyProperties[T any](target interface{}) T {
	var t T
	voValue := reflect.ValueOf(target)
	dtoValue := reflect.New(reflect.TypeOf(t)).Elem()

	for i := 0; i < voValue.NumField(); i++ {
		dtoField := dtoValue.Field(i)
		voField := voValue.Field(i)
		dtoField.Set(voField)
	}
	return dtoValue.Interface().(T)
}

func GetRuntimeStack() string {
	var buf [4096]byte
	n := runtime.Stack(buf[:], false)
	return string(buf[:n])
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
func GenericSort[T any](arr []T, less func(T, T) bool) []T {
	sort.Slice(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
	return arr
}
