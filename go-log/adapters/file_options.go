package adapters

var defaultFileOptions = fileOptions{
	Filepath: "logs/",
	MaxSize:  1 << 29,
}

type fileOptions struct {
	Filepath string
	MaxSize  int64
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
