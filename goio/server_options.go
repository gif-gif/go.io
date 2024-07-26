package goio

var defaultOptions = &options{
	noAccessPath: map[string]struct{}{
		"/favicon.ico": {},
	},
	noLogPath: map[string]struct{}{
		"/favicon.ico": {},
	},
	corsHeaders: []string{
		"Content-Type", "Content-Length",
		"Accept", "Referer", "User-Agent", "Authorization",
		"X-Requested-Id", "X-Request-Timestamp", "X-Request-Sign",
		"X-Request-AppId", "X-Request-Source", "X-Request-Token",
		"X-Client-Id", "X-Client-Token",
	},
	encryptionExcludeUris: map[string]struct{}{},
}

type options struct {
	pprofEnable bool

	serverName string
	env        Environment

	corsHeaders  []string
	noAccessPath map[string]struct{}
	noLogPath    map[string]struct{}

	encryption            *Encryption
	encryptionEnable      bool
	encryptionExcludeUris map[string]struct{}
}

type Option interface {
	apply(opts *options)
}

type funcOption struct {
	f func(opts *options)
}

func newFuncOption(f func(opts *options)) *funcOption {
	return &funcOption{f: f}
}

func (f funcOption) apply(opts *options) {
	f.f(opts)
}

// 开启分析
func PProfEnableOption(pprofEnable bool) Option {
	return newFuncOption(func(opts *options) {
		opts.pprofEnable = pprofEnable
	})
}

// 服务名称
func ServerNameOption(serverName string) Option {
	return newFuncOption(func(opts *options) {
		opts.serverName = serverName
	})
}

// 运行环境
func EnvOption(env Environment) Option {
	return newFuncOption(func(opts *options) {
		opts.env = env
	})
}

// 跨域
func CorsHeaderOption(corsHeaders ...string) Option {
	return newFuncOption(func(opts *options) {
		opts.corsHeaders = append(opts.corsHeaders, corsHeaders...)
	})
}

// 禁止访问的path
func NoAccessPathsOption(noAccessPaths ...string) Option {
	return newFuncOption(func(opts *options) {
		for _, i := range noAccessPaths {
			opts.noAccessPath[i] = struct{}{}
		}
	})
}

// 不记录日志的path
func NoLogPathsOption(noLogPaths ...string) Option {
	return newFuncOption(func(opts *options) {
		for _, i := range noLogPaths {
			opts.noLogPath[i] = struct{}{}
		}
	})
}

// 启用加密传输
func EnableEncryptionOption(encryptKey, encryptSecret string, excludeUris ...string) Option {
	return newFuncOption(func(opts *options) {
		opts.encryptionEnable = true
		opts.encryption = &Encryption{Key: encryptKey, Secret: encryptSecret}
		for _, uri := range excludeUris {
			opts.encryptionExcludeUris[uri] = struct{}{}
		}
	})
}
