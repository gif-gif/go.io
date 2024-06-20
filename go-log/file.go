package golog

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type FileAdapter struct {
	filepath string
	filename string
	fh       *os.File

	maxSize int64
	count   int

	ch chan []byte
	mu sync.Mutex
}

func NewFileLog(opts ...FileOption) *Logger {
	return New(NewFileAdapter(opts...))
}

func NewFileAdapter(opt ...FileOption) *FileAdapter {
	opts := defaultFileOptions
	for _, o := range opt {
		o.apply(&opts)
	}

	fa := &FileAdapter{
		filepath: opts.Filepath,
		maxSize:  opts.MaxSize,
		ch:       make(chan []byte, runtime.NumCPU()*2),
	}

	if l := len(fa.filepath); fa.filepath[l-1:] != "/" {
		fa.filepath += "/"
	}

	if _, err := os.Stat(fa.filepath); os.IsNotExist(err) {
		os.MkdirAll(fa.filepath, 0755)
	}

	var (
		nw  = time.Now()
		ymd = nw.Format("20060102")
	)

	files, _ := filepath.Glob(fa.filepath + ymd + "_*.log")
	fa.count = len(files)

	return fa
}

func (fa *FileAdapter) Write(msg *Message) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		c := <-fa.ch
		fa.writeHandle(c)
	}()

	fa.ch <- msg.JSON()
}

func (fa *FileAdapter) writeHandle(b []byte) {
	fa.mu.Lock()
	defer func() { fa.mu.Unlock() }()

	var (
		nw       = time.Now()
		ymd      = nw.Format("20060102")
		filename = ymd + ".log"
	)

	if filename != fa.filename {
		fa.filename = filename
		fa.closeFile()
	}

	if fa.fh == nil {
		if err := fa.openFile(); err != nil {
			return
		}
	}

	if fa.maxSize > 0 {
		if err := fa.cutFile(ymd); err != nil {
			return
		}
	}

	if fa.fh != nil {
		fa.fh.Write(b)
		fa.fh.Write([]byte("\n"))
	}
}

func (fa *FileAdapter) closeFile() {
	if fa.fh == nil {
		return
	}

	fa.fh.Close()
	fa.fh = nil
}

func (fa *FileAdapter) openFile() (err error) {
	fa.fh, err = os.OpenFile(fa.filepath+fa.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err.Error())
	}
	return
}

func (fa *FileAdapter) cutFile(ymd string) (err error) {
	var info os.FileInfo

	info, err = os.Stat(fa.filepath + fa.filename)
	if err != nil {
		log.Println(err.Error())
		return
	}
	if info.Size() < fa.maxSize {
		return
	}

	fa.count++

	filename := fmt.Sprintf("%s_%d.log", fa.filepath+ymd, fa.count)
	if err = os.Rename(fa.filepath+fa.filename, filename); err != nil {
		log.Println(err.Error())
		return
	}

	if fa.fh != nil {
		fa.fh.Close()
		fa.fh = nil
	}

	if err = fa.openFile(); err != nil {
		return
	}
	return
}
