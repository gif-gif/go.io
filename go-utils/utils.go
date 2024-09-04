package goutils

import (
	"fmt"
	"github.com/gogf/gf/util/gconv"
	"golang.org/x/crypto/bcrypt"
	"math"
	"math/rand"
	"reflect"
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

// 常用签名验证, sign Sha1 小写
func CheckSignSha1(secret, nonce string, linkSignTimeout int64, ts int64, sign string) bool {
	if linkSignTimeout == 0 {
		linkSignTimeout = 20
	}
	tsStep := time.Now().Unix() - ts
	if math.Abs(gconv.Float64(tsStep)) > gconv.Float64(linkSignTimeout) { //连接超时
		return false
	}

	args := []string{secret, gconv.String(ts), nonce}
	sort.Strings(args)
	// 小写
	serverSign := SHA1([]byte(strings.Join(args, "")))
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
