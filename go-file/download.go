package gofile

import (
	"context"
	"fmt"
	goerror "github.com/gif-gif/go.io/go-error"
	"github.com/zeromicro/go-zero/core/logx"
	"net/http"
)

type GoDownload struct {
	w        http.ResponseWriter
	r        *http.Request
	Logger   logx.Logger
	ctx      context.Context
	filename string
}

func NewGoDownload(ctx context.Context, filename string, w http.ResponseWriter, r *http.Request) *GoDownload {
	return &GoDownload{
		Logger:   logx.WithContext(ctx),
		ctx:      ctx,
		w:        w,
		r:        r,
		filename: filename,
	}
}

func WriteWithBuffer(w http.ResponseWriter, data []byte) error {
	// 使用 32KB 的缓冲区
	const bufferSize = 4 * 1024

	// 分块处理数据
	for i := 0; i < len(data); i += bufferSize {
		end := i + bufferSize
		if end > len(data) {
			end = len(data)
		}

		chunk := data[i:end]
		n, err := w.Write(chunk)
		if err != nil {
			return err
		}
		if n != len(chunk) {
			return fmt.Errorf("short write: %d/%d", n, len(chunk))
		}

		// 使用 Flush 确保数据发送
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	return nil
}

func (g *GoDownload) Write(data []byte) error {
	return WriteWithBuffer(g.w, data)
}

func (g *GoDownload) WriteString(data string) error {
	return g.Write([]byte(data))
}

func (g *GoDownload) Error(err error) {
	http.Error(g.w, "UploadError:"+err.Error(), http.StatusInternalServerError)
}

// 二进制读取输出 每次读取4096字节 输出
func (g *GoDownload) Output(filePath string) error {
	err := ReadFileChunks(filePath, 4096, func(chunk []byte) error {
		err := g.Write(chunk)
		return err
	})

	if err != nil {
		g.Error(err)
		return err
	}
	return nil
}

// 一行一行输出
func (g *GoDownload) OutputByLine(filePath string) error {
	err := ReadLines(filePath, func(chunk string) error {
		err := g.Write([]byte(chunk + "\n"))
		return err
	})
	if err != nil {
		g.Error(err)
		return err
	}
	return nil
}

func (g *GoDownload) SetFileHeaders() error {
	g.w.Header().Set("Content-Type", "application/octet-stream")
	g.w.Header().Set("Content-Disposition", "attachment; filename="+g.filename)
	g.w.Header().Set("Cache-Control", "no-cache")
	g.w.Header().Set("Connection", "keep-alive")
	g.w.WriteHeader(http.StatusOK)
	// 使用 Flush 确保数据发送
	if f, ok := g.w.(http.Flusher); ok {
		f.Flush()
	} else {
		return goerror.NewErrorMsg(uint32(http.StatusInternalServerError), "Streaming unsupported!")
	}
	return nil
}
