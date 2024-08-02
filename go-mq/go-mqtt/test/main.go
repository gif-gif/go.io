package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 定义一个彩票结构体
type Lottery struct {
	odds      map[int]float64
	pool      float64
	totalBets float64
	mu        sync.Mutex
}

// 初始化彩票结构体
func NewLottery(initialPool float64) *Lottery {
	odds := make(map[int]float64)
	for i := 0; i <= 9; i++ {
		odds[i] = 10.0 // 初始固定赔率为10倍
	}
	return &Lottery{odds: odds, pool: initialPool, totalBets: 0}
}

// 更新赔率函数，根据奖池总金额重新计算赔率
func (l *Lottery) UpdateOdds() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for i := 0; i <= 9; i++ {
		l.odds[i] = l.pool / (l.totalBets / 10.0) // 根据奖池总金额和总投注金额重新计算赔率
		if l.odds[i] < 1.0 {
			l.odds[i] = 1.0 // 确保赔率不低于1.0
		}
	}
}

// 下注函数，返回获胜金额和中奖号码
func (l *Lottery) Bet(number int, amount float64) (float64, int) {
	if number < 0 || number > 9 {
		fmt.Println("无效的数字，请选择0到9之间的数字。")
		return 0, -1
	}

	l.mu.Lock()
	l.totalBets += amount // 更新总投注金额
	l.mu.Unlock()

	// 模拟彩票抽奖
	rand.Seed(time.Now().UnixNano())
	winner := rand.Intn(10)

	fmt.Printf("中奖号码是: %d\n", winner)

	l.mu.Lock()
	defer l.mu.Unlock()
	if winner == number {
		potentialWinnings := amount * l.odds[number]
		if potentialWinnings > l.pool {
			potentialWinnings = l.pool // 确保奖励不超过奖池
		}
		l.pool -= potentialWinnings // 更新奖池
		return potentialWinnings, winner
	}
	l.pool += amount // 未中奖金额进入奖池
	return 0, winner
}

func main() {
	initialPool := 1000.0 // 初始奖池金额
	lottery := NewLottery(initialPool)

	// 模拟a投注9
	aNumber := 9
	aAmount := 100.0 // a投注100元

	// 获取a下注时的赔率
	lottery.mu.Lock()
	aOdds := lottery.odds[aNumber]
	lottery.mu.Unlock()
	fmt.Printf("a选择的数字是: %d, 赔率是: %.2f\n", aNumber, aOdds)

	aWinnings, aWinner := lottery.Bet(aNumber, aAmount)
	if aWinner == aNumber {
		fmt.Printf("恭喜a，你赢了%.2f元！\n", aWinnings)
	} else {
		fmt.Printf("很遗憾，a没有中奖。中奖号码是: %d\n", aWinner)
	}

	// 更新赔率，模拟时间推移
	lottery.UpdateOdds()

	// 模拟b投注1
	bNumber := 1
	bAmount := 50.0 // b投注50元

	// 获取b下注时的赔率
	lottery.mu.Lock()
	bOdds := lottery.odds[bNumber]
	lottery.mu.Unlock()
	fmt.Printf("b选择的数字是: %d, 赔率是: %.2f\n", bNumber, bOdds)

	bWinnings, bWinner := lottery.Bet(bNumber, bAmount)
	if bWinner == bNumber {
		fmt.Printf("恭喜b，你赢了%.2f元！\n", bWinnings)
	} else {
		fmt.Printf("很遗憾，b没有中奖。中奖号码是: %d\n", bWinner)
	}

	// 打印剩余奖池金额
	fmt.Printf("剩余奖池金额: %.2f元\n", lottery.pool)
}
