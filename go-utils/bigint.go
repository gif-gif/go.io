package goutils

import (
	"math/big"
)

// 加
func BigIntAdd(num1 string, num2 string) string {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString(num2, 10)
	x.Add(x, y)
	return x.String()
}

// 减
func BigIntReduce(num1 string, num2 string) string {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString("-"+num2, 10)
	x.Add(x, y)
	return x.String()
}

// 乘
func BigIntMul(num1 string, num2 string) string {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString(num2, 10)
	x.Mul(x, y)
	return x.String()
}

// 除
func BigIntDiv(num1 string, num2 string) string {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString(num2, 10)
	x.Div(x, y)
	return x.String()
}

// 取模
func BigIntMod(num1 string, num2 string) string {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString(num2, 10)
	x.Mod(x, y)
	return x.String()
}

// 比大小，大于返回1，等于返回0，小于返回-1
func BigIntCmp(num1 string, num2 string) int {
	x, _ := new(big.Int).SetString(num1, 10)
	y, _ := new(big.Int).SetString(num2, 10)
	return x.Cmp(y)
}
