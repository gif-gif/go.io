package gocontext

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	"os"
	"os/signal"
	"syscall"
)

func WithCancelx() *Context {
	ctx, cancel := context.WithCancel(context.TODO())
	sig := make(chan os.Signal)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				golog.Error(r)
			}
		}()

		for ch := range sig {
			switch ch {
			case syscall.SIGUSR1: // kill -USR1
			case syscall.SIGUSR2: // kill -USR2
			case syscall.SIGHUP: // kill -1
			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL: // kill -9 or ctrl+c
				cancel()
			}
		}
	}()

	return WithParent(ctx)
}
