package mw

import (
	"context"
	desc "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"

	"github.com/ulule/limiter/v3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var handlerWhitelistWithoutRateLimiter = map[string]struct{}{
	desc.OrdersService_AcceptOrder_FullMethodName:   {},
	desc.OrdersService_GetHistory_FullMethodName:    {},
	desc.OrdersService_ListOrders_FullMethodName:    {},
	desc.OrdersService_ListReturns_FullMethodName:   {},
	desc.OrdersService_ProcessOrders_FullMethodName: {},
	desc.OrdersService_ReturnOrder_FullMethodName:   {},
}

func RateLimiterInterceptor(limiter *limiter.Limiter) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if _, ok := handlerWhitelistWithoutRateLimiter[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		sender := "unknown"
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if s, ok := md["sender"]; ok {
				sender = s[0]
			}
		}

		limiterCtx, err := limiter.Get(ctx, sender)
		if err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
		if limiterCtx.Reached {
			return nil, status.Error(codes.ResourceExhausted, "rate limited")
		}

		return handler(ctx, req)
	}
}
