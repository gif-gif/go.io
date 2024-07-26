package goio

import (
	"context"
	gofile "github.com/gif-gif/go.io/go-file"
	"github.com/gin-gonic/gin"

	"google.golang.org/grpc/metadata"
	"strings"
)

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
			md.Set("trace-id", v)
		}
		if v := RequestId(c); v != "" {
			arr := gofile.Trace(2)
			if l := len(arr); l > 0 {
				md.Set("caller", strings.Join(arr, ", "))
			}
		}
	}
	return metadata.NewOutgoingContext(context.TODO(), md)
}
