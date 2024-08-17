package gocontext

import (
	"context"
	golog "github.com/gif-gif/go.io/go-log"
	"os"
	"os/signal"
	"syscall"
)

var (
	__signal chan os.Signal
	__ctx    context.Context
	__cancel context.CancelFunc
)

func init() {
	__signal = make(chan os.Signal)
	__ctx, __cancel = context.WithCancel(context.TODO())

	signal.Notify(__signal, syscall.SIGHUP, syscall.SIGUSR1, syscall.SIGUSR2,
		syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				golog.WithTag("gocontext").Error(r)
			}
		}()

		for ch := range __signal {
			switch ch {
			case syscall.SIGUSR1: // kill -USR1

			case syscall.SIGUSR2: // kill -USR2

			case syscall.SIGHUP: // kill -1

			case syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT, syscall.SIGKILL: // kill -9 or ctrl+c
				__cancel()
			}
		}
	}()
}

func Cancel() context.Context {
	return __ctx
}
