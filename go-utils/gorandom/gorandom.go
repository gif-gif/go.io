package gorandom

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// RandomGenerator 结构体定义
type RandomGenerator struct {
	rng *rand.Rand
}

// 创建新的随机数生成器
func NewRandomGenerator() *RandomGenerator {
	// Go 1.20+ 自动使用随机种子
	return &RandomGenerator{
		rng: rand.New(rand.NewPCG(1, 2)), // 使用固定种子 (1, 2)
	}
}

// 创建带自定义种子的随机数生成器
func NewRandomGeneratorWithSeed(seed uint64) *RandomGenerator {
	return &RandomGenerator{
		rng: rand.New(rand.NewPCG(seed, 0)),
	}
}

// 创建带自定义种子的随机数生成器
func NewRandomGeneratorWithSeed2(seed1 uint64, seed2 uint64) *RandomGenerator {
	return &RandomGenerator{
		rng: rand.New(rand.NewPCG(seed1, seed2)),
	}
}

// RandomGenerator 的方法
func (r *RandomGenerator) Int(min, max int) int {
	if min > max {
		min, max = max, min
	}
	return r.rng.IntN(max-min+1) + min
}

func (r *RandomGenerator) Float() float64 {
	return r.rng.Float64()
}

func (r *RandomGenerator) FloatRange(min, max float64) float64 {
	return min + r.rng.Float64()*(max-min)
}

func (r *RandomGenerator) String(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.rng.IntN(len(charset))]
	}
	return string(b)
}

func (r *RandomGenerator) Bool() bool {
	return r.rng.IntN(2) == 1
}

func (r *RandomGenerator) Choice(slice []string) string {
	if len(slice) == 0 {
		return ""
	}
	return slice[r.rng.IntN(len(slice))]
}

func (r *RandomGenerator) Shuffle(slice []int) {
	r.rng.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}

// 使用示例
func TestGoRandom() {
	// 示例 1: 基本使用
	fmt.Println("=== 示例 1: 基本使用 ===")
	rg := NewRandomGenerator()

	// 生成随机数
	fmt.Printf("随机整数 (1-100): %d\n", rg.Int(1, 100))
	fmt.Printf("随机浮点数: %.4f\n", rg.Float())
	fmt.Printf("随机字符串: %s\n", rg.String(10))
	fmt.Printf("随机布尔值: %v\n", rg.Bool())

	// 示例 2: 使用固定种子（可重现的随机序列）
	fmt.Println("\n=== 示例 2: 固定种子 ===")
	rg1 := NewRandomGeneratorWithSeed(12345)
	rg2 := NewRandomGeneratorWithSeed(12345)

	now := time.Now()
	seed1 := uint64(now.UnixNano())
	seed2 := uint64(now.UnixNano() >> 32)
	rg3 := NewRandomGeneratorWithSeed2(seed1, seed2)

	fmt.Printf("生成器固定两个种子: %d, %d, %d\n", rg3.Int(1, 100), rg3.Int(1, 100), rg3.Int(1, 100))
	// 两个生成器会产生相同的随机数序列
	fmt.Printf("生成器1: %d, %d, %d\n", rg1.Int(1, 100), rg1.Int(1, 100), rg1.Int(1, 100))
	fmt.Printf("生成器2: %d, %d, %d\n", rg2.Int(1, 100), rg2.Int(1, 100), rg2.Int(1, 100))

	// 示例 3: 游戏应用
	fmt.Println("\n=== 示例 3: 游戏应用 ===")
	game := &GameExample{
		rng: NewRandomGenerator(),
	}
	game.Play()

	// 示例 4: 多个独立的生成器
	fmt.Println("\n=== 示例 4: 多个独立生成器 ===")
	playerRng := NewRandomGenerator() // 玩家的随机数
	enemyRng := NewRandomGenerator()  // 敌人的随机数
	itemRng := NewRandomGenerator()   // 物品的随机数

	fmt.Printf("玩家伤害: %d\n", playerRng.Int(10, 20))
	fmt.Printf("敌人伤害: %d\n", enemyRng.Int(5, 15))
	fmt.Printf("掉落物品: %s\n", itemRng.Choice([]string{"剑", "盾", "药水", "金币"}))

	// 示例 5: 并发使用
	fmt.Println("\n=== 示例 5: 并发使用 ===")
	concurrentExample()
}

// 游戏示例
type GameExample struct {
	rng *RandomGenerator
}

func (g *GameExample) Play() {
	// 掷骰子
	dice := g.rng.Int(1, 6)
	fmt.Printf("掷骰子结果: %d\n", dice)

	// 暴击判定
	critChance := 0.3
	isCrit := g.rng.Float() < critChance
	fmt.Printf("是否暴击: %v\n", isCrit)

	// 随机事件
	events := []string{"遇到宝箱", "遇到怪物", "发现秘密通道", "什么都没有"}
	event := g.rng.Choice(events)
	fmt.Printf("随机事件: %s\n", event)
}

// 并发示例 - 每个 goroutine 使用自己的生成器
func concurrentExample() {
	var wg sync.WaitGroup

	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 每个 goroutine 有自己的生成器
			rng := NewRandomGeneratorWithSeed(uint64(id))

			for j := 0; j < 3; j++ {
				num := rng.Int(1, 100)
				fmt.Printf("Goroutine %d: %d\n", id, num)
			}
		}(i)
	}

	wg.Wait()
}

// 高级用法：带状态的随机数生成器
type StatefulRandom struct {
	*RandomGenerator
	callCount int
}

func NewStatefulRandom() *StatefulRandom {
	return &StatefulRandom{
		RandomGenerator: NewRandomGenerator(),
		callCount:       0,
	}
}

func (s *StatefulRandom) NextInt(min, max int) int {
	s.callCount++
	return s.Int(min, max)
}

func (s *StatefulRandom) GetCallCount() int {
	return s.callCount
}

// 测试随机数分布
func testDistribution() {
	rg := NewRandomGenerator()
	counts := make(map[int]int)

	// 生成 10000 个 1-10 的随机数
	for i := 0; i < 10000; i++ {
		num := rg.Int(1, 10)
		counts[num]++
	}

	fmt.Println("\n=== 随机数分布测试 ===")
	for i := 1; i <= 10; i++ {
		fmt.Printf("%d: %d 次 (%.2f%%)\n", i, counts[i], float64(counts[i])/100)
	}
}
