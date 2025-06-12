package adapters

import "time"

var defaultFileOptions = fileOptions{
	Filepath:         "logs/",
	MaxSize:          1 << 29,
	KeepDays:         7,              //日志保留7天
	ClearLogInterval: time.Hour * 12, //每12小清理一次
}

type fileOptions struct {
	Filepath         string
	MaxSize          int64
	KeepDays         int
	ClearLogInterval time.Duration
}

type FileOption interface {
	apply(*fileOptions)
}

type funcFileOption struct {
	f func(*fileOptions)
}

func (f *funcFileOption) apply(options *fileOptions) {
	f.f(options)
}

func newFuncFileOption(f func(*fileOptions)) *funcFileOption {
	return &funcFileOption{
		f: f,
	}
}

func FilePathOption(filepath string) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.Filepath = filepath
	})
}

func FileMaxSizeOption(maxSize int64) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.MaxSize = maxSize
	})
}

func KeepDaysOption(keepDays int) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.KeepDays = keepDays
	})
}

func ClearLogIntervalOption(clearLogInterval time.Duration) FileOption {
	return newFuncFileOption(func(options *fileOptions) {
		options.ClearLogInterval = clearLogInterval
	})
}
