package gohttpc

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/imroc/req/v3"
)

// Option 配置选项函数
type Option func(*Config)

// Config 高并发 HTTP 客户端配置
type Config struct {
	// 连接池
	MaxConnsPerHost     int           // 每个 host 的最大连接数（含 idle + active）
	MaxIdleConns        int           // 全局最大空闲连接数
	MaxIdleConnsPerHost int           // 每个 host 的最大空闲连接数
	IdleConnTimeout     time.Duration // 空闲连接回收时间

	// 超时
	DialTimeout         time.Duration // TCP 建连超时
	TLSHandshakeTimeout time.Duration // TLS 握手超时
	RequestTimeout      time.Duration // 单次请求整体超时

	// 并发控制
	MaxConcurrency int // 最大并发请求数（semaphore 大小）

	// 重试
	MaxRetries    int           // 最大重试次数
	RetryInterval time.Duration // 初始重试间隔（指数退避基数）

	// 其他
	BaseURL   string
	UserAgent string
	Debug     bool
}

// DefaultConfig 返回面向高并发场景的默认配置
func DefaultConfig() *Config {
	return &Config{
		MaxConnsPerHost:     128,
		MaxIdleConns:        512,
		MaxIdleConnsPerHost: 64,
		IdleConnTimeout:     90 * time.Second,

		DialTimeout:         5 * time.Second,
		TLSHandshakeTimeout: 5 * time.Second,
		RequestTimeout:      10 * time.Second,

		MaxConcurrency: 256,

		MaxRetries:    0,
		RetryInterval: 500 * time.Millisecond,

		UserAgent: "HighConcurrencyClient/1.0",
	}
}

// --- Option functions ---

func WithMaxConnsPerHost(n int) Option     { return func(c *Config) { c.MaxConnsPerHost = n } }
func WithMaxIdleConns(n int) Option        { return func(c *Config) { c.MaxIdleConns = n } }
func WithMaxIdleConnsPerHost(n int) Option { return func(c *Config) { c.MaxIdleConnsPerHost = n } }
func WithIdleConnTimeout(d time.Duration) Option {
	return func(c *Config) { c.IdleConnTimeout = d }
}
func WithDialTimeout(d time.Duration) Option { return func(c *Config) { c.DialTimeout = d } }
func WithTLSHandshakeTimeout(d time.Duration) Option {
	return func(c *Config) { c.TLSHandshakeTimeout = d }
}
func WithRequestTimeout(d time.Duration) Option { return func(c *Config) { c.RequestTimeout = d } }
func WithMaxConcurrency(n int) Option           { return func(c *Config) { c.MaxConcurrency = n } }
func WithMaxRetries(n int) Option               { return func(c *Config) { c.MaxRetries = n } }
func WithRetryInterval(d time.Duration) Option  { return func(c *Config) { c.RetryInterval = d } }
func WithBaseURL(url string) Option             { return func(c *Config) { c.BaseURL = url } }
func WithUserAgent(ua string) Option            { return func(c *Config) { c.UserAgent = ua } }
func WithDebug(on bool) Option                  { return func(c *Config) { c.Debug = on } }

// Client 高并发 HTTP 客户端，内置连接池管理和并发控制
type Client struct {
	client *req.Client
	sem    chan struct{} // 并发信号量
	cfg    *Config
}

// defaultClient 使用默认配置的全局客户端，通过包级函数直接使用
var defaultClient = New()

// Default 返回默认全局客户端实例
func Default() *Client {
	return defaultClient
}

// --- 包级快捷函数，直接使用默认客户端 ---

// Get 使用默认客户端发起 GET 请求
func Get(ctx context.Context, url string) (*req.Response, error) {
	return defaultClient.Get(ctx, url)
}

// Post 使用默认客户端发起 POST 请求
func Post(ctx context.Context, url string, body any) (*req.Response, error) {
	return defaultClient.Post(ctx, url, body)
}

// Do 使用默认客户端发起自定义请求
func Do(ctx context.Context, fn func(r *req.Request) (*req.Response, error)) (*req.Response, error) {
	return defaultClient.Do(ctx, fn)
}

// BatchDo 使用默认客户端并发执行一批任务
func BatchDo(ctx context.Context, tasks []Task) []Result {
	return defaultClient.BatchDo(ctx, tasks)
}

// New 创建高并发客户端
func New(opts ...Option) *Client {
	cfg := DefaultConfig()
	for _, o := range opts {
		o(cfg)
	}

	c := req.C().
		SetTimeout(cfg.RequestTimeout).
		SetCommonRetryCount(cfg.MaxRetries).
		SetCommonRetryBackoffInterval(cfg.RetryInterval, cfg.RetryInterval*16).
		SetCommonRetryCondition(func(resp *req.Response, err error) bool {
			// 仅在网络错误或 5xx/429 时重试
			if err != nil {
				return true
			}
			code := resp.StatusCode
			return code == http.StatusTooManyRequests || code >= http.StatusInternalServerError
		}).
		SetUserAgent(cfg.UserAgent)

	if cfg.BaseURL != "" {
		c.SetBaseURL(cfg.BaseURL)
	}
	if cfg.Debug {
		c.EnableDebugLog()
	}

	// 关键：调优底层 Transport 连接池，防止端口枯竭
	t := c.GetTransport()
	t.MaxConnsPerHost = cfg.MaxConnsPerHost
	t.MaxIdleConns = cfg.MaxIdleConns
	t.MaxIdleConnsPerHost = cfg.MaxIdleConnsPerHost
	t.IdleConnTimeout = cfg.IdleConnTimeout
	t.TLSHandshakeTimeout = cfg.TLSHandshakeTimeout
	t.DisableKeepAlives = false // 确保启用 keep-alive 复用连接

	return &Client{
		client: c,
		sem:    make(chan struct{}, cfg.MaxConcurrency),
		cfg:    cfg,
	}
}

// acquire / release 信号量
func (c *Client) acquire(ctx context.Context) error {
	select {
	case c.sem <- struct{}{}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) release() {
	<-c.sem
}

// Get 发起 GET 请求（受并发控制）
func (c *Client) Get(ctx context.Context, url string) (*req.Response, error) {
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()
	return c.client.R().SetContext(ctx).Get(url)
}

// Post 发起 POST 请求（受并发控制）
func (c *Client) Post(ctx context.Context, url string, body any) (*req.Response, error) {
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()
	return c.client.R().SetContext(ctx).SetBody(body).Post(url)
}

// Do 发起自定义请求（受并发控制）。回调 fn 用于配置 req.Request。
func (c *Client) Do(ctx context.Context, fn func(r *req.Request) (*req.Response, error)) (*req.Response, error) {
	if err := c.acquire(ctx); err != nil {
		return nil, err
	}
	defer c.release()
	return fn(c.client.R().SetContext(ctx))
}

// Raw 返回底层 req.Client，用于不需要并发控制的场景
func (c *Client) Raw() *req.Client {
	return c.client
}

// --- 批量并发执行 ---

// Task 代表一个待执行的请求任务
type Task struct {
	Name string                                                      // 任务标识（可选）
	Fn   func(ctx context.Context, c *Client) (*req.Response, error) // 执行函数
}

// Result 一个任务的执行结果
type Result struct {
	Name     string
	Response *req.Response
	Err      error
}

// BatchDo 并发执行一批任务，所有任务共享 Client 的并发上限。
// 返回顺序与 tasks 一致。
func (c *Client) BatchDo(ctx context.Context, tasks []Task) []Result {
	results := make([]Result, len(tasks))
	var wg sync.WaitGroup
	wg.Add(len(tasks))

	for i, task := range tasks {
		go func(idx int, t Task) {
			defer wg.Done()
			resp, err := t.Fn(ctx, c)
			results[idx] = Result{Name: t.Name, Response: resp, Err: err}
		}(i, task)
	}

	wg.Wait()
	return results
}
