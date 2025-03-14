package goserver

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"runtime"
	"strings"

	"google.golang.org/grpc/metadata"
)

func prettyFile(file string) string {
	index := strings.LastIndex(file, "/")
	if index < 0 {
		return file
	}

	index2 := strings.LastIndex(file[:index], "/")
	if index2 < 0 {
		return file[index+1:]
	}

	return file[index2+1:]
}

// 追踪信息
func trace(skip int) (arr []string) {
	arr = []string{}
	if skip == 0 {
		skip = 1
	}
	for i := skip; i < 16; i++ {
		_, file, line, _ := runtime.Caller(i)
		if file == "" {
			continue
		}
		if strings.Contains(file, ".pb.go") ||
			strings.Contains(file, "runtime/") ||
			strings.Contains(file, "src/testing") ||
			strings.Contains(file, "pkg/mod/") ||
			strings.Contains(file, "vendor/") {
			continue
		}
		arr = append(arr, fmt.Sprintf("%s %dL", prettyFile(file), line))
	}
	return
}

func GrpcContext(c *gin.Context) context.Context {
	md := metadata.New(map[string]string{})
	if c != nil {
		if v := c.GetString("__server_name"); v != "" {
			md.Set("server-name", v)
		}
		if v := c.GetString("__env"); v != "" {
			md.Set("env", v)
		}
		if v := RequestId(c); v != "" {
			md.Set("gotrace-id", v)
		}
		if v := RequestId(c); v != "" {
			arr := trace(2)
			if l := len(arr); l > 0 {
				md.Set("caller", strings.Join(arr, ", "))
			}
		}
	}
	return metadata.NewOutgoingContext(context.TODO(), md)
}
