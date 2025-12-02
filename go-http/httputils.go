package gohttp

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
)

var (
	// 全局单例客户端
	globalRestyClient *resty.Client
	once              sync.Once
)

func GetRestyClient() *resty.Client {
	once.Do(func() {
		globalRestyClient = NewOptimizedClient()
	})
	return globalRestyClient
}

// 创建优化的 HTTP 客户端
func NewOptimizedClient() *resty.Client {
	client := resty.New()

	// 1. 连接池优化
	client.SetTransport(&http.Transport{
		MaxIdleConns:        200,              // 最大空闲连接数
		MaxIdleConnsPerHost: 200,              // 每个host的最大空闲连接
		MaxConnsPerHost:     200,              // 每个host的最大连接数
		IdleConnTimeout:     90 * time.Second, // 空闲连接超时
		DisableKeepAlives:   false,            // 启用Keep-Alive
		DisableCompression:  false,            // 启用压缩
	})

	// 2. 超时设置
	client.SetTimeout(10 * time.Second)
	client.SetRetryCount(3)                     // 重试次数
	client.SetRetryWaitTime(1 * time.Second)    // 重试等待时间
	client.SetRetryMaxWaitTime(5 * time.Second) // 最大重试等待时间

	// 3. 重试条件
	client.AddRetryCondition(func(r *resty.Response, err error) bool {
		// 网络错误或5xx错误时重试
		return err != nil || r.StatusCode() >= 500
	})

	// 4. 通用请求头
	client.SetHeaders(map[string]string{
		"User-Agent": "Optimized-Resty-Client/1.0",
		"Accept":     "application/json",
	})

	return client
}
