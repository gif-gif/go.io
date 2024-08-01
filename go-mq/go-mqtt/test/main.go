package main

import (
	"fmt"
)

// 计算赔率
func calculateOdds(totalNumbers int) float64 {
	return float64(totalNumbers)
}

// 计算期望值
func calculateExpectedValue(totalNumbers int, prizeAmount float64, ticketCost float64) float64 {
	winningProbability := 1.0 / float64(totalNumbers)
	expectedValue := (winningProbability * prizeAmount) - ticketCost
	return expectedValue
}

func main() {
	totalNumbers := 10   // 从0到9，共10个数字
	prizeAmount := 100.0 // 奖金金额
	ticketCost := 2.0    // 每次投注金额

	odds := calculateOdds(totalNumbers)
	expectedValue := calculateExpectedValue(totalNumbers, prizeAmount, ticketCost)

	fmt.Printf("赔率: 1:%.0f\n", odds)
	fmt.Printf("期望值: %.2f\n", expectedValue)
}
