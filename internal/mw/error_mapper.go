package mw

import (
	"context"
	"errors"
	businesserror "gitlab.ozon.dev/gleb_evseev/homework/internal/domain/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MapBusinessErrorToGRPC конвертирует бизнес-ошибку в gRPC статус.
func MapBusinessErrorToGRPC(err error) error {
	if err == nil {
		return nil
	}

	var be *businesserror.Error
	if errors.As(err, &be) {
		switch be.Code {
		case businesserror.ErrorCodeNotFound:
			return status.Errorf(codes.NotFound, "%s", be.Message)
		case businesserror.ErrorCodeAlreadyExists:
			return status.Errorf(codes.AlreadyExists, "%s", be.Message)
		case businesserror.ErrorCodeExpires:
			return status.Errorf(codes.FailedPrecondition, "%s", be.Message)
		case businesserror.ErrorCodeValidation:
			return status.Errorf(codes.InvalidArgument, "%s", be.Message)
		case businesserror.ErrorCodePackage, businesserror.ErrorCodeWeight:
			return status.Errorf(codes.InvalidArgument, "%s", be.Message)
		case businesserror.ErrorCodeInternal:
			return status.Errorf(codes.Internal, "%s", be.Message)
		default:
			return status.Errorf(codes.Unknown, "%s", be.Message)

		}
	}

	return status.Errorf(codes.Internal, "%s", err.Error())

}

func ErrorMappingInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			return nil, MapBusinessErrorToGRPC(err)
		}
		return resp, nil
	}
}
