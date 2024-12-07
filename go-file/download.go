package gofile

import (
	"context"
	"fmt"
	goerror "github.com/gif-gif/go.io/go-error"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type Belt struct {
	DataChan  chan []byte
	CloseChan chan interface{}
}

type GoDownload struct {
	w        http.ResponseWriter
	r        *http.Request
	Logger   logx.Logger
	Belt     *Belt
	exitChan chan interface{}
	ctx      context.Context
	filename string
}

func NewGoDownload(ctx context.Context, filename string, w http.ResponseWriter, r *http.Request) *GoDownload {
	return &GoDownload{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		w:      w,
		r:      r,
		Belt: &Belt{
			DataChan:  make(chan []byte, 20),
			CloseChan: make(chan interface{}),
		},
		exitChan: make(chan interface{}),
		filename: filename,
	}
}

func (g *GoDownload) WaitDone() {
	<-g.exitChan
	close(g.exitChan)
}

// 输出文件完成后记得关闭，否则会一直等待导致内存泄漏
func (g *GoDownload) Close() {
	g.Belt.CloseChan <- 1
}

func (g *GoDownload) Write(data []byte) {
	g.Belt.DataChan <- data
}

func (g *GoDownload) Error(err error) {
	http.Error(g.w, "UploadError:"+err.Error(), http.StatusInternalServerError)
}

// 二进制读取输出 每次读取4096字节 输出
func (g *GoDownload) Output(filePath string) error {
	err := ReadFileChunks(filePath, 4096, func(chunk []byte) error {
		g.Write(chunk)
		return nil
	})

	g.Close()
	if err != nil {
		g.Error(err)
		return err
	}
	return nil
}

// 一行一行输出
func (g *GoDownload) OutputByLine(filePath string) error {
	err := ReadLines(filePath, func(chunk string) error {
		g.Write([]byte(chunk + "\n"))
		return nil
	})

	g.Close()
	if err != nil {
		g.Error(err)
		return err
	}
	return nil
}

func (g *GoDownload) Run() {
	defer func() {
		if err := recover(); err != nil {
			g.Logger.Errorf("GoDownload Run recover error:%+v,msg:%s", err, goerror.GetStack())
		}
	}()

	flusher, ok := g.w.(http.Flusher)
	if !ok {
		http.Error(g.w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	g.w.Header().Set("Content-Type", "application/octet-stream")
	g.w.Header().Set("Content-Disposition", "attachment; filename="+g.filename)
	g.w.Header().Set("Cache-Control", "no-cache")
	g.w.Header().Set("Connection", "keep-alive")

	g.w.WriteHeader(http.StatusOK)
	flusher.Flush()

ForLoop:
	for {
		select {
		case data := <-g.Belt.DataChan:
			if len(data) > 0 {
				fmt.Fprintf(g.w, "%s", data)
			}
			flusher.Flush()
		case <-g.Belt.CloseChan:
			close(g.Belt.CloseChan)
			g.exitChan <- 1
			close(g.Belt.DataChan) //????
			break ForLoop
		}
	}
}
