package main

import (
	"context"
	"fmt"
	"time"

	"github.com/gif-gif/go.io/go-http/gohttpc"
	"github.com/imroc/req/v3"
)

func main() {
	ctx := context.Background()

	// --- 示例 1: 直接用包级函数，零配置 ---
	resp, err := gohttpc.Get(ctx, "https://httpbin.org/get")
	if err != nil {
		fmt.Println("GET error:", err)
	} else {
		fmt.Println("GET status:", resp.StatusCode)
	}

	// --- 示例 2: 包级 POST ---
	resp, err = gohttpc.Post(ctx, "https://httpbin.org/post", map[string]string{
		"hello": "world",
	})
	if err != nil {
		fmt.Println("POST error:", err)
	} else {
		fmt.Println("POST status:", resp.StatusCode)
	}

	// --- 示例 3: 包级 Do 自定义请求 ---
	resp, err = gohttpc.Do(ctx, func(r *req.Request) (*req.Response, error) {
		return r.
			SetHeader("X-Custom", "value").
			SetQueryParam("page", "1").
			Get("https://httpbin.org/get")
	})
	if err != nil {
		fmt.Println("Do error:", err)
	} else {
		fmt.Println("Do status:", resp.StatusCode)
	}

	// --- 示例 4: 需要自定义配置时才 New ---
	client := gohttpc.New(
		gohttpc.WithMaxConcurrency(100),
		gohttpc.WithMaxConnsPerHost(50),
		gohttpc.WithRequestTimeout(10*time.Second),
		gohttpc.WithMaxRetries(2),
	)

	urls := []string{
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/1",
		"https://httpbin.org/delay/1",
	}

	tasks := make([]gohttpc.Task, len(urls))
	for i, u := range urls {
		tasks[i] = gohttpc.Task{
			Name: fmt.Sprintf("req-%d", i),
			Fn: func(ctx context.Context, c *gohttpc.Client) (*req.Response, error) {
				return c.Get(ctx, u)
			},
		}
	}

	start := time.Now()
	results := client.BatchDo(ctx, tasks)
	fmt.Printf("\nBatch: %d tasks done in %v\n", len(results), time.Since(start))
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("  %s: error=%v\n", r.Name, r.Err)
		} else {
			fmt.Printf("  %s: status=%d\n", r.Name, r.Response.StatusCode)
		}
	}
}
