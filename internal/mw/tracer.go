package mw

import (
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// TraceInterceptor — gRPC interceptor, добавляющий трейсинги от запросов.
func TraceInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {

	span := trace.SpanFromContext(ctx)
	if span != nil && span.IsRecording() {
		span.SetAttributes(attribute.String("rpc.method", info.FullMethod))
	}

	return handler(ctx, req)
}
