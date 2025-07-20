package mw

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// LoggingInterceptor — gRPC interceptor, логирующий входящие запросы, метаданные и время обработки.
func LoggingInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Printf("SendMessage: %v", req)

	start := time.Now()
	resp, err = handler(ctx, req)

	md, _ := metadata.FromIncomingContext(ctx)
	log.Printf("Method: %s, Duration: %s, Metadata: %v, Error: %v",
		info.FullMethod,
		time.Since(start),
		md,
		err,
	)

	return resp, err
}
