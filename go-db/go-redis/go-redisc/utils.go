package goredisc

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type LatencyStats struct {
	Samples map[string][]time.Duration `json:"-"` // 原始样本数据不序列化
	Total   int64                      `json:"total"`
	Avg     time.Duration              `json:"avg"`
	mu      sync.RWMutex               `json:"-"` // 互斥锁不序列化
}

// CommandStats 单个命令的统计信息
type CommandStats struct {
	Command string        `json:"command"`
	Count   int           `json:"count"`
	Average time.Duration `json:"average"`
	Min     time.Duration `json:"min"`
	Max     time.Duration `json:"max"`
	P50     time.Duration `json:"p50"`
	P90     time.Duration `json:"p90"`
	P95     time.Duration `json:"p95"`
	P99     time.Duration `json:"p99"`
	P999    time.Duration `json:"p999"`
}

// LatencyStatsJSON JSON 序列化结构体
type LatencyStatsJSON struct {
	TotalSamples int64          `json:"total_samples"`
	OverallAvg   time.Duration  `json:"overall_avg"`
	Commands     []CommandStats `json:"commands"`
	Timestamp    time.Time      `json:"timestamp"`
}

// LatencyStatsJSONString JSON 序列化结构体（字符串格式的时间）
type LatencyStatsJSONString struct {
	TotalSamples int64                `json:"total_samples"`
	OverallAvg   string               `json:"overall_avg"`
	Commands     []CommandStatsString `json:"commands"`
	Timestamp    string               `json:"timestamp"`
}

// CommandStatsString 单个命令的统计信息（字符串格式的时间）
type CommandStatsString struct {
	Command string `json:"command"`
	Count   int    `json:"count"`
	Average string `json:"average"`
	Min     string `json:"min"`
	Max     string `json:"max"`
	P50     string `json:"p50"`
	P90     string `json:"p90"`
	P95     string `json:"p95"`
	P99     string `json:"p99"`
	P999    string `json:"p999"`
}

// NewLatencyStats 创建新的延迟统计实例
func NewLatencyStats() *LatencyStats {
	return &LatencyStats{
		Samples: make(map[string][]time.Duration),
	}
}

// AddSample 添加延迟样本
func (ls *LatencyStats) AddSample(command string, latency time.Duration) {
	ls.mu.Lock()
	defer ls.mu.Unlock()

	if ls.Samples[command] == nil {
		ls.Samples[command] = make([]time.Duration, 0)
	}

	ls.Samples[command] = append(ls.Samples[command], latency)
	ls.Total++

	// 重新计算平均值
	ls.calculateAverage()
}

// calculateAverage 计算总体平均延迟
func (ls *LatencyStats) calculateAverage() {
	var totalDuration time.Duration
	var count int64

	for _, samples := range ls.Samples {
		for _, sample := range samples {
			totalDuration += sample
			count++
		}
	}

	if count > 0 {
		ls.Avg = totalDuration / time.Duration(count)
	}
}

// GetCommandPercentile 获取指定命令的百分位延迟
func (ls *LatencyStats) GetCommandPercentile(command string, percentile float64) time.Duration {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	samples := ls.Samples[command]
	if len(samples) == 0 {
		return 0
	}

	sorted := make([]time.Duration, len(samples))
	copy(sorted, samples)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	index := int(float64(len(sorted)) * percentile / 100.0)
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// GetCommandAverage 获取指定命令的平均延迟
func (ls *LatencyStats) GetCommandAverage(command string) time.Duration {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	samples := ls.Samples[command]
	if len(samples) == 0 {
		return 0
	}

	var total time.Duration
	for _, sample := range samples {
		total += sample
	}

	return total / time.Duration(len(samples))
}

// getCommandMinMax 获取指定命令的最小值和最大值
func (ls *LatencyStats) getCommandMinMax(command string) (time.Duration, time.Duration) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	samples := ls.Samples[command]
	if len(samples) == 0 {
		return 0, 0
	}

	min, max := samples[0], samples[0]
	for _, sample := range samples {
		if sample < min {
			min = sample
		}
		if sample > max {
			max = sample
		}
	}

	return min, max
}

// GetCommandCount 获取指定命令的样本数量
func (ls *LatencyStats) GetCommandCount(command string) int {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	return len(ls.Samples[command])
}

// GetAllCommands 获取所有测试的命令列表
func (ls *LatencyStats) GetAllCommands() []string {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	commands := make([]string, 0, len(ls.Samples))
	for cmd := range ls.Samples {
		commands = append(commands, cmd)
	}

	return commands
}

// ToJSON 将延迟统计转换为 JSON 字节数组（time.Duration 格式）
func (ls *LatencyStats) ToJSON() (LatencyStatsJSON, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	jsonData := LatencyStatsJSON{
		TotalSamples: ls.Total,
		OverallAvg:   ls.Avg,
		Commands:     make([]CommandStats, 0),
		Timestamp:    time.Now(),
	}

	for cmd := range ls.Samples {
		if len(ls.Samples[cmd]) == 0 {
			continue
		}

		min, max := ls.getCommandMinMax(cmd)

		cmdStats := CommandStats{
			Command: cmd,
			Count:   ls.GetCommandCount(cmd),
			Average: ls.GetCommandAverage(cmd),
			Min:     min,
			Max:     max,
			P50:     ls.GetCommandPercentile(cmd, 50),
			P90:     ls.GetCommandPercentile(cmd, 90),
			P95:     ls.GetCommandPercentile(cmd, 95),
			P99:     ls.GetCommandPercentile(cmd, 99),
			P999:    ls.GetCommandPercentile(cmd, 99.9),
		}

		jsonData.Commands = append(jsonData.Commands, cmdStats)
	}

	return jsonData, nil
}

// ToJSONString 将延迟统计转换为 JSON 字节数组（字符串格式的时间，便于阅读）
func (ls *LatencyStats) ToJSONString() (LatencyStatsJSONString, error) {
	ls.mu.RLock()
	defer ls.mu.RUnlock()

	jsonData := LatencyStatsJSONString{
		TotalSamples: ls.Total,
		OverallAvg:   ls.Avg.String(),
		Commands:     make([]CommandStatsString, 0),
		Timestamp:    time.Now().Format("2006-01-02 15:04:05"),
	}

	for cmd := range ls.Samples {
		if len(ls.Samples[cmd]) == 0 {
			continue
		}

		min, max := ls.getCommandMinMax(cmd)

		cmdStats := CommandStatsString{
			Command: cmd,
			Count:   ls.GetCommandCount(cmd),
			Average: ls.GetCommandAverage(cmd).String(),
			Min:     min.String(),
			Max:     max.String(),
			P50:     ls.GetCommandPercentile(cmd, 50).String(),
			P90:     ls.GetCommandPercentile(cmd, 90).String(),
			P95:     ls.GetCommandPercentile(cmd, 95).String(),
			P99:     ls.GetCommandPercentile(cmd, 99).String(),
			P999:    ls.GetCommandPercentile(cmd, 99.9).String(),
		}

		jsonData.Commands = append(jsonData.Commands, cmdStats)
	}

	return jsonData, nil
}

// benchmarkRedisLatency 使用新的结构进行延迟基准测试
func BenchmarkRedisLatency(rdb *redis.ClusterClient, duration time.Duration) *LatencyStats {
	ctx := context.Background()
	stats := NewLatencyStats()

	fmt.Printf("Running latency benchmark for %v...\n", duration)

	start := time.Now()

	// 定义要测试的命令
	commands := map[string]func() error{
		"PING": func() error {
			return rdb.Ping(ctx).Err()
		},
		"SET": func() error {
			return rdb.Set(ctx, fmt.Sprintf("test_key_%d", time.Now().UnixNano()), "test_value", 20*time.Second).Err()
		},
		"GET": func() error {
			return rdb.Get(ctx, "test_key").Err()
		},
		"INCR": func() error {
			return rdb.Incr(ctx, "test_counter").Err()
		},
		"LPUSH": func() error {
			return rdb.LPush(ctx, "test_list", "test_value").Err()
		},
	}

	commandNames := make([]string, 0, len(commands))
	for cmd := range commands {
		commandNames = append(commandNames, cmd)
	}

	cmdIndex := 0
	for time.Since(start) < duration {
		// 轮询执行不同命令
		cmdName := commandNames[cmdIndex%len(commandNames)]
		cmdFunc := commands[cmdName]

		cmdStart := time.Now()
		err := cmdFunc()
		latency := time.Since(cmdStart)

		if err == nil {
			stats.AddSample(cmdName, latency)
		}

		cmdIndex++

		// 控制请求频率
		time.Sleep(10 * time.Millisecond)
	}

	return stats
}
