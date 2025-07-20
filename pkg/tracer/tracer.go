// Package tracer предоставляет инициализацию и работу с OpenTelemetry трассировкой,
//
// Основные возможности:
//   - Инициализация трассировщика с OTLP-экспортером (InitTracer)
//   - Создание новых спанов (Start)
//
// Использование:
//
//	shutdown := tracer.InitTracer(ctx, "app-name", otlpEndpoint)
//	defer shutdown()
//
//	ctx, span := tracer.Start(ctx, "operation")
//	defer span.End()
package tracer

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	ottrace "go.opentelemetry.io/otel/trace"
)

var globalTracer = otel.Tracer("pvz-service")

// InitTracer — инициализирует OpenTelemetry трассировку с OTLP-экспортером по gRPC.
func InitTracer(ctx context.Context, serviceName, Endpoint string) func() {
	exp, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(Endpoint),
		otlptracegrpc.WithDialOption(),
	)
	if err != nil {
		log.Fatalf("failed to create OTLP exporter: %v", err)
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(semconv.ServiceName(serviceName)),
	)
	if err != nil {
		log.Fatalf("failed to create resource: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exp),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return func() {
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown tracer: %v", err)
		}
	}

}

func Start(ctx context.Context, spanName string) (context.Context, ottrace.Span) {
	return globalTracer.Start(ctx, spanName)
}
