package grpc

import (
	"context"
	"errors"
	businesserror "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"gitlab.ozon.dev/gleb_evseev/homework/internal/domain/pvz"
	pkg "gitlab.ozon.dev/gleb_evseev/homework/pkg/pvz"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type OrdersHandler interface {
	Accept(ctx context.Context, order pvz.AcceptOrderInput) error
	Return(ctx context.Context, orderID pvz.ReturnOrderInput) error
	Process(ctx context.Context, orders pvz.ProcessOrdersInput) (pvz.ProcessResult, error)
	ListOrders(ctx context.Context, opts pvz.ListOrdersInput) (*pvz.OrdersList, error)
	ListReturns(ctx context.Context, opts pvz.ListReturnsInput) (*pvz.ReturnsList, error)
	OrderHistory(ctx context.Context) (*pvz.OrderHistoryList, error)
}
type OrdersServer struct {
	pkg.UnimplementedOrdersServiceServer
	ordersHandler OrdersHandler
}

func NewOrdersServer(handler OrdersHandler) *OrdersServer {
	return &OrdersServer{
		ordersHandler: handler,
	}
}

// MapBusinessErrorToGRPC конвертирует бизнес-ошибку в gRPC статус.
func MapBusinessErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var be *businesserror.Error
	if errors.As(err, &be) {
		switch be.Code {
		case businesserror.ErrorCodeNotFound:
			return status.Errorf(codes.NotFound, be.Message)
		case businesserror.ErrorCodeAlreadyExists:
			return status.Errorf(codes.AlreadyExists, be.Message)
		case businesserror.ErrorCodeExpires:
			return status.Errorf(codes.FailedPrecondition, be.Message)
		case businesserror.ErrorCodeValidation:
			return status.Errorf(codes.InvalidArgument, be.Message)
		case businesserror.ErrorCodePackage, businesserror.ErrorCodeWeight:
			return status.Errorf(codes.InvalidArgument, be.Message)
		case businesserror.ErrorCodeInternal:
			return status.Errorf(codes.Internal, be.Message)
		default:
			return status.Errorf(codes.Unknown, be.Message)
		}
	}

	return status.Errorf(codes.Internal, err.Error())
}
