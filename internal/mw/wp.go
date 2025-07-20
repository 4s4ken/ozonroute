package mw

import (
	"context"
	"gitlab.ozon.dev/gleb_evseev/homework/pkg/workerpool"
	"google.golang.org/grpc"
)

// WorkerPoolInterceptor — gRPC interceptor, внедряющий WorkerPool в контекст запроса.
func WorkerPoolInterceptor(wp workerpool.WorkerPool) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		ctxWithWP := workerpool.WorkerPoolValue(ctx, wp)
		return handler(ctxWithWP, req)
	}
}
