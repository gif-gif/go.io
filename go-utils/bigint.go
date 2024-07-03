package goutils

import (
	"github.com/shopspring/decimal"
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

func SubFloat64(a float64, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	_c, _ := _a.Sub(_b).Float64()
	return _c
}

func MulFloat64(a float64, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	_c, _ := _a.Mul(_b).Float64()
	return _c
}

func AddFloat64(a float64, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	_c, _ := _a.Add(_b).Float64()
	return _c
}

func DivFloat64(a float64, b float64) float64 {
	_a := decimal.NewFromFloat(a)
	_b := decimal.NewFromFloat(b)
	_c, _ := _a.Div(_b).Float64()
	return _c
}
