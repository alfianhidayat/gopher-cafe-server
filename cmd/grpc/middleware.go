package main

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

func TimeoutMiddleware() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		timeout, cancelFunc := context.WithTimeout(ctx, time.Second*2)
		defer cancelFunc()
		return handler(timeout, req)
	}
}
