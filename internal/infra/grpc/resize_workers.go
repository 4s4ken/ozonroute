package grpc

import (
	"context"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ResizeWorkerPool — обрабатывает gRPC-запрос на изменения количества воркеров в workerpool.
func (s *OrdersServer) ResizeWorkerPool(ctx context.Context, req *pkg.ResizeWorkerPoolRequest) (*pkg.ResizeWorkerPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "request is nil")
	}
	if req.Workers <= 0 {
		return nil, status.Error(codes.InvalidArgument, "workers must be positive")
	}

	err := s.ordersHandler.ResizeWorkerPool(ctx, int(req.Workers))
	if err != nil {
		return &pkg.ResizeWorkerPoolResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pkg.ResizeWorkerPoolResponse{
		Success: true,
		Message: "worker pool resized successfully",
	}, nil
}
