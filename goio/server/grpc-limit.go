package server

import (
	"context"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/syncx"
	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

// Grpc 调用限流 limit 是最大并发数
func GrpcLimit(s *zrpc.RpcServer, limit int) {
	l := syncx.NewLimit(limit)
	s.AddUnaryInterceptors(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		if l.TryBorrow() {
			defer func() {
				if err := l.Return(); err != nil {
					logx.Error(err)
				}
			}()
			return handler(ctx, req)
		} else {
			logx.Errorf("concurrent connections over %d, rejected with code %d",
				limit, http.StatusServiceUnavailable)
			return nil, status.Error(codes.Unavailable, "concurrent connections over limit")
		}
	})
}
