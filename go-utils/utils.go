package goutils

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"github.com/samber/lo"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"reflect"
	"regexp"
	"runtime"
	"slices"
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
func GenericSort[T any](arr []T, less func(T, T) bool) {
	sort.Slice(arr, func(i, j int) bool {
		return less(arr[i], arr[j])
	})
}

// 把缺失的数字填充到数组中
func FillMissingNumbers(nums []int64, max int64) []int64 {
	// 创建一个新的切片来存储结果
	var result []int64
	// 从 1 开始
	current := int64(1)

	// 遍历给定的数字
	for _, num := range nums {
		// 填充中间缺失的数字
		for current < num {
			result = append(result, current)
			current++
		}
		// 添加当前数字
		result = append(result, num)
		current = num + 1 // 更新当前数字到下一个
	}

	// 如果还有剩余的数字，继续填充
	for i := current; i <= max; i++ { // 假设我们想填充到 20
		result = append(result, i)
	}

	return result
}

func GetPageCount(total int64, pageSize int64) (totalPages int64) {
	return int64(math.Ceil(float64(total) / float64(pageSize)))
}

// 下一页
func AfterPage(page int64, pageCount int64) int64 {
	if page <= 0 {
		page = 1
	}
	after := page + 1
	if after > pageCount {
		after = -1
	}
	return after
}

func BeforePage(page int64) int64 {
	before := page - 1
	if before <= 0 {
		before = -1
	}
	return before
}

// 是不是数字
func IsNumeric(str string) bool {
	re := regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)
	// 使用正则表达式匹配字符串
	ok := re.MatchString(str)
	return ok
}

func IsInt(str string) bool {
	re := regexp.MustCompile("^[0-9]+$")
	// 使用正则表达式匹配字符串
	ok := re.MatchString(str)
	return ok
}

func Sum(list []int) int {
	return lo.Sum(list)
}
