package gorandom

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand/v2"
	"time"
)

// 1. 生成随机整数（指定范围）
func RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return rand.IntN(max-min+1) + min
}

// 2. 生成随机浮点数（0.0-1.0）
func RandomFloat() float64 {
	return rand.Float64()
}

// 3. 生成指定范围的随机浮点数
func RandomFloatRange(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

// 4. 生成随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.IntN(len(charset))]
	}
	return string(b)
}

// 5. 生成安全的随机字符串（使用 crypto/rand）
func SecureRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, err := cryptorand.Int(cryptorand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[n.Int64()]
	}
	return string(b), nil
}

// 6. 生成随机布尔值
func RandomBool() bool {
	return rand.IntN(2) == 1
}

// 7. 从切片中随机选择一个元素
func RandomChoice[T any](slice []T) T {
	if len(slice) == 0 {
		var zero T
		return zero
	}
	return slice[rand.IntN(len(slice))]
}

// 8. 打乱切片顺序（洗牌）
func Shuffle[T any](slice []T) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// 9. 生成随机密码
func RandomPassword(length int, includeSpecial bool) string {
	charset := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if includeSpecial {
		charset += "!@#$%^&*()_+-=[]{}|;:,.<>?"
	}

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.IntN(len(charset))]
	}
	return string(password)
}

// 10. 生成随机十六进制字符串
func RandomHex(length int) string {
	const hex = "0123456789abcdef"
	b := make([]byte, length)
	for i := range b {
		b[i] = hex[rand.IntN(16)]
	}
	return string(b)
}

// 12. 按权重随机选择
func WeightedRandomChoice[T any](items []T, weights []int) T {
	if len(items) == 0 || len(items) != len(weights) {
		var zero T
		return zero
	}

	totalWeight := 0
	for _, w := range weights {
		totalWeight += w
	}

	r := rand.IntN(totalWeight)
	for i, w := range weights {
		r -= w
		if r < 0 {
			return items[i]
		}
	}

	return items[len(items)-1]
}

// 13. 生成正态分布随机数
func RandomNormal(mean, stdDev float64) float64 {
	return rand.NormFloat64()*stdDev + mean
}

// 14. 生成随机颜色（RGB）
func RandomColor() string {
	return fmt.Sprintf("#%02x%02x%02x",
		rand.IntN(256),
		rand.IntN(256),
		rand.IntN(256))
}

// 15. 生成随机日期时间
func RandomDateTime(start, end time.Time) time.Time {
	delta := end.Unix() - start.Unix()
	sec := rand.Int64N(delta) + start.Unix()
	return time.Unix(sec, 0)
}

// 示例用法
func Test() {
	// 注意：Go 1.20+ 不再需要手动设置种子
	// 旧版本需要：rand.Seed(time.Now().UnixNano())

	// 1. 随机整数
	fmt.Printf("Random int (1-100): %d\n", RandomInt(1, 100))

	// 2. 随机浮点数
	fmt.Printf("Random float: %.4f\n", RandomFloat())

	// 3. 随机字符串
	fmt.Printf("Random string: %s\n", RandomString(10))

	// 4. 安全随机字符串
	secureStr, _ := SecureRandomString(16)
	fmt.Printf("Secure random string: %s\n", secureStr)

	// 5. 随机选择
	fruits := []string{"apple", "banana", "orange", "grape"}
	fmt.Printf("Random fruit: %s\n", RandomChoice(fruits))

	// 6. 洗牌
	numbers := []int{1, 2, 3, 4, 5}
	Shuffle(numbers)
	fmt.Printf("Shuffled numbers: %v\n", numbers)

	// 7. 随机密码
	fmt.Printf("Random password: %s\n", RandomPassword(12, true))

	// 9. 加权随机
	items := []string{"common", "rare", "epic", "legendary"}
	weights := []int{70, 20, 8, 2}
	fmt.Printf("Weighted random: %s\n", WeightedRandomChoice(items, weights))

	// 10. 随机颜色
	fmt.Printf("Random color: %s\n", RandomColor())

	// 11. 随机日期
	start := time.Now().AddDate(-1, 0, 0)
	end := time.Now()
	fmt.Printf("Random date: %s\n", RandomDateTime(start, end).Format("2006-01-02"))
}
