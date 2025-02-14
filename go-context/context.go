package gocontext

import (
	"context"
	"encoding/json"
	golog "github.com/gif-gif/go.io/go-log"
)

// 应用程序全局上下文
var (
	__ctx context.Context
)

func init() {
	__ctx = WithCancel().Context
}

func Cancel() context.Context {
	return __ctx
}

// 上下文
type Context struct {
	context.Context
	Log *golog.Entry
	v   map[string]any
}

func (ctx *Context) WithLog() *Context {
	ctx.Log = golog.Default().WithTag()
	return ctx
}

func (ctx *Context) WithParent(parent context.Context) *Context {
	ctx.Context = parent
	return ctx
}

func (ctx *Context) WithValue(key string, value any) *Context {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	ctx.v[key] = value
	return ctx
}

func (ctx *Context) Value(key string) any {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	if v, ok := ctx.v[key]; ok {
		return v
	}
	return nil
}

func (ctx *Context) Values() map[string]any {
	if ctx.v == nil {
		ctx.v = map[string]any{}
	}
	return ctx.v
}

func (ctx *Context) Json() []byte {
	v := ctx.Values()
	b, _ := json.Marshal(&v)
	return b
}

func (ctx *Context) String() string {
	return string(ctx.Json())
}

func WithCancel() *Context {
	ctx, _ := context.WithCancel(context.TODO())
	//sig := make(chan os.Signal)
	//syscall.SIGUSR1, syscall.SIGUSR2,
	//signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGBREAK, syscall.SIGKILL)
	//
	//go func() {
	//	defer func() {
	//		if r := recover(); r != nil {
	//			golog.Error(r)
	//		}
	//	}()
	//
	//	for ch := range sig {
	//		switch ch {
	//		case syscall.SIGUSR1: // kill -USR1
	//		case syscall.SIGUSR2: // kill -USR2
	//		case syscall.SIGHUP: // kill -1
	//		case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL: // kill -9 or ctrl+c
	//			cancel()
	//		}
	//	}
	//}()

	return WithParent(ctx)
}

func WithLog() *Context {
	return Default().WithLog()
}

func WithParent(parent context.Context) *Context {
	return Default().WithParent(parent)
}

func Default() *Context {
	return &Context{}
}
